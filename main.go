package main

import (
	"context"
	"github.com/idawud/medium-monitor/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

var (
	l = log.New(os.Stdout, "medium-monitor ", log.LstdFlags)
)

func main() {

	// Setting up our only route
	sm := mux.NewRouter()
	sm.HandleFunc("/feedback", WebsocketEndpoint)


	// server connection tuning
	server := &http.Server{
		Addr: ":8080",
		Handler: sm,
		IdleTimeout:120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout:1*time.Second,
	}

	log.Println("Server running on http://localhost:8080/")
	// Run server on port :8080
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Graceful shutdown with Ctrl+C
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = server.Shutdown(ctx)
}

// Setup the socket connection buffer size for read & write
// & allowing CORS on all origins so that we'll test it with our client
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true },
}

// the /feedback route implementation
func WebsocketEndpoint(writer http.ResponseWriter, request *http.Request) {
	l.Println("Main WebSocket (Start)")
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		l.Println("Error: ", err)
	}
	l.Println("Client successfully connected")

	// read message from client and write back
	readerAndWriter(ws)
}

// Listens for message on socket connection and also sens back messages to the client
func readerAndWriter(ws *websocket.Conn) {
	// run infinitely
	for  {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			l.Println(err)
			return
		}
		// read message on received
		l.Println("Client Message: ", string(p))

		// write message to client every 15secs
		for  {
			// get the server status
			availability, err := service.GetAllAvailability()
			if err != nil {
				l.Println("Error: ", err)
			} else {
				// write the message to the client which connected to it
				if err := ws.WriteMessage(messageType, availability); err != nil {
					l.Println("Error: ", err)
					_ = ws.Close()
				} else {
					l.Println("New Data Published")
				}
			}

			time.Sleep(time.Second * 15)
		}

	}
}
