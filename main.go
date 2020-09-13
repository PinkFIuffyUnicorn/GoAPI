// Package classification of GO API
//
// Documentation for GO API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define structure for User
// swagger:response user
type user struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
}

// Define structure for Group
type group struct {
	Name string `json:"Name"`
}

// Connect to MongoDB
func mongoDbConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://admin:admin@goapi.qzo53.mongodb.net/UsersAndGroups?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// swagger:route POST /users users addUser
// Adds a User record to the Users collection
// responses:
//	200: user
func addUser(w http.ResponseWriter, r *http.Request) {
	var user user
	reqBody, _ := ioutil.ReadAll(r.Body)
	reqBody = []byte(strings.ToLower(string(reqBody)))
	json.Unmarshal(reqBody, &user)

	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	usersColletion := usersAndGroupsDatabase.Collection("Users")
	usersResult, err := usersColletion.InsertOne(context.Background(), bson.D{
		{Key: "Name", Value: user.Name},
		{Key: "Password", Value: user.Password},
		{Key: "Email", Value: user.Email},
	})

	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(usersResult)
}

// swagger:route GET /users users getUser
// Returns a User record
// responses:
//	200: user
func getUser(w http.ResponseWriter, r *http.Request) {
	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	usersColletion := usersAndGroupsDatabase.Collection("Users")

	params := r.URL.Query()
	comibnedString := "{"

	for param, paramValue := range params {
		comibnedString = comibnedString + "\"" + param + "\"" + ":" + paramValue[0] + ","
	}

	if comibnedString != "{" {
		comibnedString = strings.TrimSuffix(comibnedString, ",")
	}
	comibnedString = comibnedString + "}"

	var bsonMap bson.M
	err := json.Unmarshal([]byte(comibnedString), &bsonMap)
	if err != nil {
		log.Fatal(err)
	}

	filterUsers, err := usersColletion.Find(context.Background(), bsonMap)
	if err != nil {
		log.Fatal(err)
	}

	var userFiltered []bson.M

	if err = filterUsers.All(context.Background(), &userFiltered); err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(userFiltered)
}

// swagger:route PUT /users users updateUser
// Updates a User record in the Users collection
// responses:
//	200: user
func updateUser(w http.ResponseWriter, r *http.Request) {
	var user user
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)

	userName := user.Name
	userPassword := user.Password
	userEmail := user.Email

	dict := make(map[string]string)
	if userName != "" {
		dict["Name"] = userName
	}
	if userPassword != "" {
		dict["Password"] = userPassword
	}
	if userEmail != "" {
		dict["Email"] = user.Email
	}

	vars := mux.Vars(r)
	idRequest, _ := vars["id"]
	id, _ := primitive.ObjectIDFromHex(idRequest)

	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	usersColletion := usersAndGroupsDatabase.Collection("Users")

	counter := 0

	for key, value := range dict {
		counter++
		result, err := usersColletion.UpdateOne(
			context.Background(),
			bson.M{"_id": id},
			bson.D{
				{"$set", bson.D{{key, value}}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(bson.M{"_id": idRequest, "FieldsUpdated": counter})
}

// Delete a User record
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRequest, _ := vars["id"]
	id, _ := primitive.ObjectIDFromHex(idRequest)

	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	usersColletion := usersAndGroupsDatabase.Collection("Users")

	result, err := usersColletion.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(bson.M{"ID": idRequest, "deletedCount": result.DeletedCount})
}

// Add a Group record
func addGroup(w http.ResponseWriter, r *http.Request) {
	var group group
	reqBody, _ := ioutil.ReadAll(r.Body)
	reqBody = []byte(strings.ToLower(string(reqBody)))
	json.Unmarshal(reqBody, &group)

	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	groupsColletion := usersAndGroupsDatabase.Collection("Groups")

	cursor, err := groupsColletion.Find(context.Background(), bson.M{"Name": group.Name})
	if err != nil {
		log.Fatal(err)
	}

	if cursor.RemainingBatchLength() == 0 {
		groupsResult, err := groupsColletion.InsertOne(context.Background(), bson.D{
			{Key: "Name", Value: group.Name},
		})

		if err != nil {
			log.Fatal(err)
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "	")
		enc.Encode(groupsResult)
	} else {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "	")
		enc.Encode(bson.M{"response": "Group Name already exists"})
	}
}

// Delete a Group record
func deleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRequest, _ := vars["id"]
	id, _ := primitive.ObjectIDFromHex(idRequest)

	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	groupsColletion := usersAndGroupsDatabase.Collection("Groups")

	result, err := groupsColletion.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(bson.M{"_id": idRequest, "deletedCount": result.DeletedCount})
}

// Get a Group record
func getGroup(w http.ResponseWriter, r *http.Request) {
	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	groupsColletion := usersAndGroupsDatabase.Collection("Groups")

	params := r.URL.Query()
	comibnedString := "{"

	for param, paramValue := range params {
		comibnedString = comibnedString + "\"" + strings.Title(strings.ToLower(param)) + "\"" + ":" + paramValue[0] + ","
	}

	if comibnedString != "{" {
		comibnedString = strings.TrimSuffix(comibnedString, ",")
	}
	comibnedString = comibnedString + "}"

	var bsonMap bson.M
	err := json.Unmarshal([]byte(comibnedString), &bsonMap)
	if err != nil {
		log.Fatal(err)
	}

	filterGroups, err := groupsColletion.Find(context.Background(), bsonMap)
	if err != nil {
		log.Fatal(err)
	}

	var groupFiltered []bson.M

	if err = filterGroups.All(context.Background(), &groupFiltered); err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(groupFiltered)
}

// Update a Group record
func updateGroup(w http.ResponseWriter, r *http.Request) {
	var group group
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &group)

	groupName := group.Name
	if groupName != "" {
		vars := mux.Vars(r)
		idRequest, _ := vars["id"]
		id, _ := primitive.ObjectIDFromHex(idRequest)

		client := mongoDbConnect()
		defer client.Disconnect(context.Background())

		usersAndGroupsDatabase := client.Database("UsersAndGroups")
		groupsColletion := usersAndGroupsDatabase.Collection("Groups")

		result, err := groupsColletion.UpdateOne(
			context.Background(),
			bson.M{"_id": id},
			bson.D{
				{"$set", bson.D{{"Name", group.Name}}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(result)

		enc := json.NewEncoder(w)
		enc.SetIndent("", "	")
		enc.Encode(bson.M{"_id": idRequest, "updatedCount": result.ModifiedCount})
	}
}

// Define URL requests and run the app
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", addUser).Methods("POST")
	myRouter.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/users", getUser).Methods("GET")
	myRouter.HandleFunc("/groups", getGroup).Methods("GET")
	myRouter.HandleFunc("/groups", addGroup).Methods("POST")
	myRouter.HandleFunc("/groups/{id}", updateGroup).Methods("PUT")
	myRouter.HandleFunc("/groups/{id}", deleteGroup).Methods("DELETE")

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	myRouter.Handle("/docs", sh)
	myRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
