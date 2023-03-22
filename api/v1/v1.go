package v1

import (
	"fmt"
	"net/http"
	"strings"
)

// Incoming() -- responds to v1 requests that made it past front.go
func Incoming(w http.ResponseWriter, r *http.Request) {
	// remove leading forward slash.
	path := strings.TrimPrefix(r.URL.Path, "/")

	// split path into array to grab the value we need
	pathSplit := strings.Split(path, "/")
	apiCall := pathSplit[2]

	// now we'll start procesing the request
	// from here check the apiCall value and route appropriately
	switch apiCall {
	case "customers":
		fmt.Fprintf(w, "api/v1/customers\n")
		customers(w, r)
	case "employees":
		fmt.Fprintf(w, "api/v1/employees\n")
		folders(w, r)
	}
}
