package main

import (
	"fmt"
	"log"
	"net/http"
	"context"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Element struct {
	ID primitive.ObjectID `bson:"_id"`
	Content string `bson:"text"`
}

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/add", add).Methods("POST")
	log.Fatal(http.ListenAndServe(":3001", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("simpleCrudGoMongo").Collection("elements")
}

func add(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var element Element
	element.ID = primitive.NewObjectID()
	json.Unmarshal(reqBody, &element)
	
	_, err := collection.InsertOne(ctx, element)
	if err != nil {
		log.Fatal(err)
	}
}	