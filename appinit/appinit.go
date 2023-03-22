package appinit

import (
	"os"
	"restapi/api"
	"restapi/logging"

	"github.com/joho/godotenv"
)

// var Run bool = true //while true ap will run Global value

var (
	Run string = "run" // Global flag we can use to shut down the app, gracefully or not.
	// mongodbReady   bool
	// ProjectID      string
	// Bucket       = os.Getenv("BUCKET_NAME")
	// gacJson string //GOOGLE_APPLICATION_CREDENTIALS json file location
	// topics  []string
)

// Init() -- this is for initializing the application
// we hard code some values here for now but they should
// be moved off to a config file or database.
// This is were we would set up our clients for dependencies. DB, Kafka, Redis, etc.
func Init() {
	// Set Timezone
	// TODO: Why am I setting this?
	os.Setenv("TZ", "America/New_York")

	// initialize logging system
	logger := logging.GetLoggerInstance()
	logger.SetLogLevel(1) // this should be a flag we can set while the app is running. DB, config.file. Something that can be updated via the api.
	logger.SetLogFileNames("./logs/activity.log", "./logs/error.log", "./logs/service.log")
	logger.ServiceLogging("appinit.Init", "Logging system initialized")

	// Inititialize environment variables
	if err := godotenv.Load(); err != nil {
		logger.ServiceLogging("appinit.Init.godotenv", "No .env file found. Will use local variables")
	} else {
		logger.ServiceLogging("appinit.Init.godotenv", ".env file found. Using variables stored here")
	}

	// // initialize external dependencies here
	// Kafka/pubsub/redis/etc/SQL

	// and finally we'll start the REST api
	// we do this last so the end point isn't
	// available until all other services are up.
	logger.ServiceLogging("appinit.Init", "Registering and initializing API endpoint")
	api.RegisterControllers()

	// Log that the service is properly initialized
	logger.ServiceLogging("appinit.Init", "Application is properly initialized")

}
