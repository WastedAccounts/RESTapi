package api

import (
	"fmt"
	"log"
	"net/http"

	v1 "restapi/api/v1"
	"strings"

	"github.com/gorilla/mux"
)

// RegisterControllers() -- sets up the api endpoint and begins listening
func RegisterControllers() {
	// Startup mux router
	router := mux.NewRouter().StrictSlash(true)
	// first we'll redirect anything we don't know how or want to handle to a 404 message
	router.NotFoundHandler = http.HandlerFunc(nopath)
	// now we'll take in everything from the /api path
	router.PathPrefix("/api").HandlerFunc(api)
	// here we'll start listening in the background for api calls
	go func() {
		log.Fatal(http.ListenAndServe(":3000", router))
	}()
}

// nopath() -- Anything we don't know or want we send to 404
func nopath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404 not found")
}

// api() -- Here we can extend to later version of the api
func api(w http.ResponseWriter, r *http.Request) {
	apiVersion := getApiVersion(w, r)
	if apiVersion == "v1" {
		v1.Incoming(w, r)
	} else {
		// if not found we'll redirect to 404
		nopath(w, r)
	}
}

// getApiVersion() -- checks the path to make sure it contains enough values to pass
// if not we direct it to nopath() to return a 404
// pulls the api version from the path and returns it to api()
func getApiVersion(w http.ResponseWriter, r *http.Request) string {

	// /////////// Start troubleshooting code to see what we have coming in on the request
	// // react to the incoming host url
	// // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// ///////////
	// //// used for troubleshooting how to handle incoming requests.
	// //// Print incoming request
	// fmt.Println("header-referer:", r.Header.Get("Referer"))
	// fmt.Println("referer-referer:", r.Referer())
	// fmt.Println("r.URL.EscapedPath():", r.URL.EscapedPath())
	// fmt.Println("r.URL.Query():", r.URL.Query())
	// // // Create return string
	// var request []string
	// // Add the request string
	// url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	// request = append(request, url)
	// // Add the host
	// request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// // Loop through headers
	// for name, headers := range r.Header {
	// 	name = strings.ToLower(name)
	// 	for _, h := range headers {
	// 		request = append(request, fmt.Sprintf("%v: %v", name, h))
	// 	}
	// }
	// // If this is a POST, add post data
	// if r.Method == "POST" {
	// 	r.ParseForm()
	// 	request = append(request, "\n")
	// 	request = append(request, r.Form.Encode())
	// }
	// // Return the request as a string
	// fmt.Println(strings.Join(request, "\n"))
	// /////////// End troubleshooting code

	//This cuts off the leading forward slash.
	path := strings.TrimPrefix(r.URL.Path, "/")

	// split path into array to grab the value I need
	pathSplit := strings.Split(path, "/")
	apiVersion := pathSplit[1]

	// if the path doesn't contain a value after the version we nopath() it and end the call
	if len(pathSplit) < 3 {
		apiVersion = "nopath"
	}

	return apiVersion

}

// // encodeResponseAsJSON - Not in use
// func encodeResponseAsJSON(data interface{}, w io.Writer) {
// 	enc := json.NewEncoder(w)
// 	enc.Encode(data)
// 	//return data
// }
