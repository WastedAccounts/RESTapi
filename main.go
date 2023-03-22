package main

import (
	"log"
	"restapi/appinit"
	"restapi/logging"
	"time"
)

func main() {

	// initialize app
	log.Println("SERVICE: Initializing Server")
	appinit.Init()

	logging.GetLoggerInstance().ServiceLogging("main", "Initialization Complete")

	// start pinging logs to report service is up -- dev feature, remove in prod
	upTime := 0
	if appinit.Run == "run" {
		// flag app to shutdown -- appinit.Run = false
		for appinit.Run == "run" {
			// this should be removed or extended. right now it
			// just signals the app is up while debugging.
			log.Println("Monitoring folders -- Server Uptime:", upTime)
			time.Sleep(1 * time.Minute)
			upTime++
		}
	} else if appinit.Run != "run" {
		// we can extend failure scenarios here
		// graceful shutdown or fatal messages
		// can be discerned and logged
		if appinit.Run == "failedstartup" {
			// failed started up messaging
			logging.GetLoggerInstance().ErrorLogging("main", "Failed to startup")
			logging.GetLoggerInstance().ErrorLogging("main", "Please review the logs")
		}
	}
	log.Println("VERBOSE: Service is shutting down")

}
