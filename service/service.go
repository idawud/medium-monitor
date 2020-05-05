package service

import (
"encoding/json"
"net/http"
"time"
)

// endpoint on the server to test against
var ENDPOINTS = []string{
	"https://idawud.tech/myserice/api/v1/todos",
	"http://localhost:8080/api/v2/todos",
}

// Make a [HEAD] request on any point on the service server
// if the status code isn't between 200 & 299 then something
// went wrong on the server
// Else everything is ok
func CheckEndpointAvailable(url string) bool {
	res, err := http.Head(url)
	if err != nil {
		return false
	}
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return true
	}
	return false
}

// check the availability of all the endpoints and return the them as
// map of the key=url,value=availability true/false & addition
// for the timestamp of reading
func GetAllAvailability() ([]byte, error) {
	var result =  make(map[string]interface{})
	for _, ep := range ENDPOINTS{
		result[ep] = CheckEndpointAvailable(ep)
	}
	result["timestamp"] = time.Now().String()
	return json.Marshal(result)
}
