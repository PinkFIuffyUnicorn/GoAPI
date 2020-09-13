package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
}

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

func getUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	client := mongoDbConnect()
	defer client.Disconnect(context.Background())

	usersAndGroupsDatabase := client.Database("UsersAndGroups")
	usersColletion := usersAndGroupsDatabase.Collection("Users")

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

func updateUser(w http.ResponseWriter, r *http.Request) {
	var user user
	reqBody, _ := ioutil.ReadAll(r.Body)
	//reqBody = []byte(strings.ToLower(string(reqBody)))
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

	for key, value := range dict {
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
		log.Println(result)
	}

	// enc := json.NewEncoder(w)
	// enc.SetIndent("", "	")
	// enc.Encode("asdf")
}

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

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", addUser).Methods("POST")
	myRouter.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/users", getUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
