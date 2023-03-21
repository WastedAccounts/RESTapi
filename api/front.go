package api

import (
	"fmt"
	"log"
	v1 "mft/monitor/api/v1"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func RegisterControllers() {
	// Startup mux router
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(nopath)
	router.PathPrefix("/api/v1").HandlerFunc(apiV1)
	log.Fatal(http.ListenAndServe(":8080", router))

}

// handle incoming request from the router
func nopath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404 not found")
}

func apiV1(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("apiV1body", r.Body)
	check, apiCall := checkRequestPath(w, r)
	if check == true {
		if apiCall == "customers" {
			fmt.Fprintf(w, "api/v1/customers\n")
			v1.CustomerCalls(w, r)
		} else if apiCall == "folders" {
			fmt.Fprintf(w, "api/v1/folders\n")
			v1.FolderCalls(w, r)
		}
	} else {
		fmt.Fprintf(w, "404 not found\n")
	}
}

func checkRequestPath(w http.ResponseWriter, r *http.Request) (bool, string) {
	check := false

	/////////// Start troubleshooting code
	// react to the incoming host url
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	///////////
	//// used for troubleshooting how to handle incoming requests.
	//// Print incoming request
	fmt.Println("header-referer:", r.Header.Get("Referer"))
	fmt.Println("referer-referer:", r.Referer())
	fmt.Println("r.URL:", r.URL.EscapedPath())
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
	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	fmt.Println(strings.Join(request, "\n"))
	/////////// End troubleshooting code

	//This cuts off the leading forward slash.
	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// split path into array to grab the value I need
	pathSplit := strings.Split(path, "/")
	apiCall := pathSplit[2]

	// create array of acceptable path values for the api
	vals := []string{"customers", "folders"}

	// check if the incoming path is in that
	check = contains(vals, apiCall)
	// fmt.Println("check", check)

	// return a result
	return check, apiCall

}

// contains - checks if a string exists in a string array
// -- pass in a string array and the string you're looking for in that array
// -- returns a true or false
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// //========================================

// api := newApiController()
// http.Handle("/api", *api)
// http.Handle("/api/", *api)

// // loginController - handles values after .com/
// type apiController struct {
// 	apiPattern *regexp.Regexp
// }

// // newLoginController - entry point from front.go
// func newApiController() *apiController {
// 	return &apiController{
// 		apiPattern: regexp.MustCompile(`^/api/(\d+)/?`),
// 	}
// }

// func (ac apiController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path == "/api" {
// 		switch r.Method {
// 		case http.MethodGet:
// 			w.WriteHeader(http.StatusAccepted)
// 			// submit := r.FormValue("submit")
// 			// if submit == "users" {
// 			// 	// pageLoadUsers(w)
// 			// } else if submit == "movements" {
// 			// 	// pageLoadMovements(w, r)
// 			// } else if submit == "workouts" {
// 			// 	// pageLoadWorkouts()
// 			// } else if submit == "user" {
// 			// 	// pageLoadUser(w, r)
// 			// } else {
// 			// 	// pageLoadAdmin(w)
// 			// }
// 		case http.MethodPost:
// 			// submit := r.FormValue("submit")
// 			// err := r.ParseForm()
// 			// changeValue := map[string]bool{
// 			// 	"Activate":   true,
// 			// 	"Deactivate": true,
// 			// 	"User":       true,
// 			// 	"Moderator":  true,
// 			// 	"Admin":      true,
// 			// }
// 			// if err != nil {
// 			// 	log.Fatalf("Failed to decode postFormByteSlice: %v", err)
// 			// }
// 			//movements := r.FormValue("movements")

// 		default:
// 			w.WriteHeader(http.StatusNotImplemented)
// 		}
// 	}
// }

// // encodeResponseAsJSON - Not is use not sure what it did
// func encodeResponseAsJSON(data interface{}, w io.Writer) {
// 	enc := json.NewEncoder(w)
// 	enc.Encode(data)
// 	//return data
// }
