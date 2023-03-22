package v1

import (
	"fmt"
	"net/http"
	"restapi/logging"
	"strings"
)

func logRequests(w http.ResponseWriter, r *http.Request) {

	// // // Now we'll pull some data from the request for logging
	// // Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// Return the request and write to logs
	logging.GetLoggerInstance().ActivityLogging("api.v1.folders", fmt.Sprint(strings.Join(request, "\n")))
}
