package logging

import (
	"fmt"
	"log"
	"mft/ruleprocessor/cloudsql"
)

// error log insert statement
var loggingErrorExec = "INSERT INTO error_logs (service,log_content) VALUES ('%s','%s');"

// ErrorLogging() -- writes errors to error_logs table
func ErrorLogging(s1 string, err error) {
	log.Printf("ERROR: %s -- %s", s1, err)
	cloudsql.InsertErrorRecord(fmt.Sprintf(loggingErrorExec, s1, err))
}

// activity log insert statement
var loggingActivityExec = "INSERT INTO activity_logs (service,log_content) VALUES ('%s','%s');"

// ActivityLogging() -- writes activity to activity_logs table
func ActivityLogging(s1 string, s2 string) error {
	// TODO: add debug logging to console output
	// Control in config flag in database so it can be enabled while running.
	// If debug == true {write log to console} else {nothing}
	// for now we'll output to console so we can demo the service
	log.Printf("VERBOSE: %s -- %s", s1, s2)
	err := cloudsql.InsertActivityRecord(fmt.Sprintf(loggingActivityExec, s1, s2))
	if err != nil {
		ErrorLogging("logging", err)
		return err
	}
	return nil
}

// activity log insert statement
var loggingServiceExec = "INSERT INTO service_logs (package,log_content) VALUES ('%s','%s');"

// ServiceLogging() -- Verbose Logging about the system, startup, initialization, etc.
func ServiceLogging(s1 string, s2 string) {
	// s1 is package name
	// s2 is the message
	log.Printf("DEBUG: %s -- %s", s1, s2)
	err := cloudsql.InsertServiceRecord(fmt.Sprintf(loggingServiceExec, s1, s2))
	if err != nil {
		ErrorLogging("logging", err)
	}
}
