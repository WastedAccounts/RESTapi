package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
	// "mft/ruleprocessor/cloudsql"
)

type Logger struct {
	level       int
	debug       *log.Logger
	activity    *log.Logger
	service     *log.Logger
	errorLog    string
	activityLog string
	serviceLog  string
}

var once sync.Once
var logger *Logger

// GetLoggerInstance() -- instantiates our logger and ensure only one is created
func GetLoggerInstance() *Logger {
	once.Do(
		func() {
			fmt.Println("Creating Logger instance.")
			logger = &Logger{}
			//logger.ActivityLogging("logging.GetLoggerInstance", "Creating Logger instance.")
		})
	return logger
}

// ErrorLogging() -- For logging errors
func (l *Logger) ErrorLogging(pkg string, msg string) {
	if l.level <= 2 {
		f, err := os.OpenFile(l.errorLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open error.log file")
		}
		l.debug = log.New(f, "[DEBUG] -- ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
		l.debug.Printf("-- %s -- %s", pkg, msg)
	}
}

// ActivityLogging() -- for logging app activity
func (l *Logger) ActivityLogging(pkg string, msg string) {
	// this will always log application activity.
	if l.level <= 2 {
		f, err := os.OpenFile(l.activityLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open activity.log file")
		}
		l.activity = log.New(f, "[ACTIVITY] -- ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
		l.activity.Printf("-- %s -- %s", pkg, msg)
	}
}

// ServiceLogging() -- For logging service states
func (l *Logger) ServiceLogging(pkg string, msg string) {
	if l.level <= 2 {
		f, err := os.OpenFile(l.serviceLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open service.log file")
		}
		l.service = log.New(f, "[SERVICE] -- ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
		l.service.Printf("-- %s -- %s", pkg, msg)
	}
}

// SetLogLevel() -- for setting our logger's log level.
func (l *Logger) SetLogLevel(level int) {
	// check logging flag for logging level
	// Only ERROR, DEBUG, ALL implemented currently
	// ALL = 1 -- DEBUG = 2 -- ERROR = 5
	// returning 1 so all logging runs right now -- set in appinit.go
	l.level = level

	// TODO: implement a logging flag that can be switched without restarting the application
	// Config file or DB value. Create API endpoint for enabling and changing log levels
}

// SetLogFileNames() -- for setting log file names
func (l *Logger) SetLogFileNames(actlog string, errlog string, svclog string) {
	l.activityLog = actlog
	l.errorLog = errlog
	l.serviceLog = svclog
}
