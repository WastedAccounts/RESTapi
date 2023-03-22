package v1

import (
	"fmt"
	"net/http"
	"strings"
)

func Incoming(w http.ResponseWriter, r *http.Request) {
	// remove leading forward slash.
	path := strings.TrimPrefix(r.URL.Path, "/")

	// split path into array to grab the value we need
	pathSplit := strings.Split(path, "/")
	apiCall := pathSplit[2]

	// now we'll start procesing the request

	if r.Method == "GET" {
		fmt.Fprintf(w, "GET\n")
		fmt.Println("GET")
		if apiCall == "customers" {
			fmt.Fprintf(w, "api/v1/customers\n")
			fmt.Fprintf(w, "parameters: %s", r.URL.Query())
			// v1.CustomerCalls(w, r)
		} else if apiCall == "folders" {
			fmt.Fprintf(w, "api/v1/folders\n")
			fmt.Fprintf(w, "parameters: %s", r.URL.Query())
			// v1.FolderCalls(w, r)
		}
	}
	if r.Method == "POST" {
		fmt.Fprintf(w, "POST\n")
		fmt.Println("POST")
	}
	if r.Method == "PUT" {
		fmt.Fprintf(w, "PUT\n")
		fmt.Println("PUT")
	}
	if r.Method == "DELETE" {
		fmt.Fprintf(w, "DELETE\n")
		fmt.Println("DELETE")

	}
	if r.Method == "PATCH" {
		fmt.Fprintf(w, "PATCH\n")
		fmt.Println("PATCH")
	}

}
