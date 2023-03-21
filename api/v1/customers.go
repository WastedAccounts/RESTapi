package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mft/monitor/mongodb"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CustomerFiles struct
// -- used by GetSingleCustomer and GetAllCustomer structs
// -- loads files sub documents
type customerFilesStruct struct {
	Filename    string `json:"filename"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

// CustomerCalls -- Handles incoming api calls from the /api/v1/customers endpoint
func CustomerCalls(w http.ResponseWriter, r *http.Request) {
	endpoints := []string{"", "addcustomer", "getcustomer"}
	var quit bool

	//This cuts off the leading forward slash.
	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	// split path into array to grab the value I need
	pathSplit := strings.Split(path, "/")
	endpoint := ""
	if len(pathSplit) > 3 {
		endpoint = pathSplit[3]
	}
	if stringInArray(endpoints, endpoint) == false {
		stre := "501 not implemented"
		bytee := []byte(stre)
		w.Write(bytee)
		quit = true
	}
	if quit == false {
		switch r.Method {
		case http.MethodGet:
			// get all customers - TODO: need to add paging in so we don't overload the response
			if endpoint == "" {
				customers := getAllCustomers(w, r)
				cstmrs, _ := json.Marshal(customers)
				w.Write(cstmrs)
			}
			// get single customer by ARNumber
			if endpoint == "getcustomer" {
				arnumber := pathSplit[4]
				customer := getSingleCustomer(w, r, arnumber)
				cstmr, _ := json.Marshal(customer)
				w.Write(cstmr)
			}
		case http.MethodPost:
			if endpoint == "addcustomer" {
				addCustomer(w, r)
			}
		case http.MethodDelete:
		case http.MethodPut:
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// Start Get All Customers API Call //
// getAllCustomers struct
// -- used for storing customer name, ARNumber and a count of files from customers collection
type getAllCustomersStruct struct {
	Name     string `bson:"name"`
	ARnumber string `bson:"arnumber"`
	Files    int    `bson:"numberOfFiles"`
}

// getAllCustomers -- returns a list of all customers from customers collection
// -- Name, ARNumber, countOfFiles
func getAllCustomers(w http.ResponseWriter, r *http.Request) []*getAllCustomersStruct {
	var rslt getAllCustomersStruct
	var result []*getAllCustomersStruct
	conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("customers")
	filter := bson.D{{
		Key: "$project",
		Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "arnumber", Value: 1},
			{Key: "_id", Value: 0},
			{Key: "numberOfFiles", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{{Key: "$isArray", Value: "$files"}}},
					{Key: "then", Value: bson.D{{Key: "$size", Value: "$files"}}},
					{Key: "else", Value: "NA"},
				},
				}},
			},
		},
	}}

	results, err := conn.Aggregate(context.TODO(), mongo.Pipeline{filter})
	if err != nil {
		log.Fatal("errOr:", err)
	}
	for results.Next(context.TODO()) {
		// log.Println("cursor:", results)
		results.Decode(&rslt)
		result = append(result, &rslt) // (*Customers)(&result))
	}
	return result
}

// End Get All Customers API Call //

// Start Get Single Customer API Call //
// getSingleCustomer struct
// -- used for storing customer name, ARNumber and a count of files from customers collection
type getSingleCustomerStruct struct {
	Name     string `bson:"name"`
	ARnumber string `bson:"arnumber"`
	// Files   int    `bson:"numberOfFiles"`
	Filelist []customerFilesStruct `bson:"files"`
}

// getSingleCustomer -- returns a list of all customers from customers collection
// -- Name, ARNumber, countOfFiles
// TODO: this should return a list of files and not just a count
func getSingleCustomer(w http.ResponseWriter, r *http.Request, arnumber string) []*getSingleCustomerStruct {
	var rslt getSingleCustomerStruct
	var result []*getSingleCustomerStruct
	conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("customers")
	match := bson.D{{Key: "$match", Value: bson.D{{Key: "arnumber", Value: arnumber}}}}
	filter := bson.D{{
		Key: "$project",
		Value: bson.D{
			{Key: "_id", Value: 0},
		},
	}}
	results, err := conn.Aggregate(context.TODO(), mongo.Pipeline{match, filter})
	if err != nil {
		log.Fatal("errOr:", err)
	}
	for results.Next(context.TODO()) {
		// log.Println("cursor:", results)
		results.Decode(&rslt)
		result = append(result, &rslt) // (*Customers)(&result))
	}
	return result
}

// End Get Single Customer API Call //

// Start Add Customer API Call //
// addCustomer structs -- used for adding customer through the api
type addCustomerStruct struct {
	ARNumber string                `json:"arnumber"`
	Name     string                `json:"name"`
	Files    []customerFilesStruct `json:"files"`
}

// addCustomer -- addes new customer to customers collection
// -- adds customer name, arnumber, and files they need monitored
func addCustomer(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	p := addCustomerStruct{}
	err = json.Unmarshal([]byte(b), &p)
	if err != nil {
		// if error is not nil
		fmt.Println("addCustomer unmarshallError:", err)
	}
	// fmt.Println("PersonStruct: ", p)

	conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("customers")
	_, err = conn.InsertOne(context.TODO(), p)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("Customer already exists in database", err)
			stre := "Customer already exists in database\nUpdated any new folders"
			bytee := []byte(stre)
			w.Write(bytee)
		} else {
			log.Println("addCustomer error inserting into database:", err)
		}
	}
	log.Printf("%s added to the customers database", p.ARNumber)

	cf := []customerFilesStruct{}
	for _, s := range p.Files {
		cf = append(cf, s)
	}
	checkMonitoredFolders(cf)
}

// End Add Customer API Call //

// Start utility functions //
// checkMonitoredFolders -- check if the customers folders exist in folder collection
// -- if not it will add them
func checkMonitoredFolders(cf []customerFilesStruct) {
	folders := getFolders()
	if folders == nil {
		log.Println("No folders found, adding new folders")
		// // var empty *Paths
		// folders.Paths = ""
		for _, s := range cf {
			fmt.Println("s.Origin", s.Origin)
			conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("folders")
			_, err := conn.UpdateOne(context.TODO(), bson.D{{
				Key:   "_id",
				Value: folders.Id,
			}}, mongo.Pipeline{{{Key: "$set", Value: bson.D{{Key: "folders", Value: bson.D{{Key: "path", Value: s.Origin}}}}}}})
			if err != nil {
				log.Println("checkMonitored Folders error inserting into database:", err)
			}
			checkMonitoredFolders(cf)
		}
	} else {
		for _, s := range cf {
			fmt.Println("s.Origin", s.Origin)
			exists := stringInArray(folders.Paths, string(s.Origin))
			fmt.Println("folders.Paths", folders.Paths)
			fmt.Println("exists :", exists)
			if exists == false {
				// Add to DB
				// log.Printf("Path '%s' does not exist in database", string(s.Origin))
				conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("folders")
				filter := bson.D{{Key: "_id", Value: folders.Id}}
				update := bson.D{{Key: "$push", Value: bson.D{{Key: "folders", Value: s.Origin}}}}
				opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
				var updatedDoc bson.D
				err := conn.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedDoc)
				if err != nil {
					log.Println("checkMonitored Folders error inserting into database:", err)
				}
			} else {
				log.Printf("Path '%s' already exists in database", string(s.Origin))
			}
			folders = getFolders()
		}
	}
}

// structs for getFolders
type pathsStruct struct {
	// TODO: pull in ID so we can ref it instead of fighting with this
	Id    primitive.ObjectID `bson:"_id"`
	Paths []string           `bson:"folders"`
}

// getFolders -- pulls list of monitored folders
// -- used to compare exiting folders with new customer folders
func getFolders() *pathsStruct {
	var result *pathsStruct
	conn := mongodb.MongoConnPool.Database(mongodb.MongoDBName).Collection("folders")
	filter := bson.D{{}}
	collection := bson.D{{Key: "_id", Value: 1}, {Key: "folders", Value: 1}} //Key: "folders", Value: 1
	opts := options.FindOne().SetProjection(collection)

	err := conn.FindOne(context.TODO(), filter, opts).Decode(&result)
	fmt.Println("results:", result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("No folders exist yet:", err)
			result = nil
			return result
		} else {
			log.Fatal("getFolders ERROR:", err)
		}
	}
	return result
}

// stringInArray - checks if a string exists in a string array
// -- pass in a string array and the string you're looking for in that array
// -- returns a true or false
func stringInArray(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// End utility functions //
