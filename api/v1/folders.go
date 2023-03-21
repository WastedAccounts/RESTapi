package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mft/monitor/mongodb"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FolderCalls
// -- Handles incoming api calls from the /api/v1/folders endpoint
func FolderCalls(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)

	switch r.Method {
	case http.MethodGet:
		folders := getAllFolders(w, r)
		fldrs, _ := json.Marshal(folders.Folders)
		w.Write(fldrs)
	case http.MethodPost:
	case http.MethodDelete:
	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Folders Struct
// -- used for storing folders list from folders collection
type Folders struct {
	Folders []string `bson:"folders"`
}

// getAllFolders
// -- return a list of all monitored folders from the folders collection
func getAllFolders(w http.ResponseWriter, r *http.Request) *Folders {
	var result *Folders
	conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("folders")
	filter := bson.D{{}}
	collection := bson.D{{Key: "folders", Value: 1}}
	opts := options.FindOne().SetProjection(collection)

	err := conn.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// if r.URL.Path == "/api" {
// 	switch r.Method {
// 	case http.MethodGet:
// 	case http.MethodPost:
// 	case http.MethodDelete:
// 	case http.MethodPut:
// 	default:
// 		w.WriteHeader(http.StatusNotImplemented)
// 	}
// }
