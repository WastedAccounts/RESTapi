package v1

import (
	"encoding/json"
	"net/http"
)

func customers(w http.ResponseWriter, r *http.Request) {

	// handle customer GET calls
	switch r.Method {
	case "GET":
		// we can process the request and return a json response to the client
		json.NewEncoder(w).Encode(r.URL.Query())

		// and log request details to our activity log
		logRequests(w, r)
	case "POST":

	}
}
