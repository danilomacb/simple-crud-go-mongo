package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", add).Methods("POST")
	myRouter.HandleFunc("/", list)
	myRouter.Use(mux.CORSMethodMiddleware(myRouter))
	log.Fatal(http.ListenAndServe(":3001", myRouter))
}

type Element struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Content string             `bson:"content" json:"content"`
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func list(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	var elements []Element

	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var e Element
		err := cur.Decode(&e)
		if err != nil {
			log.Fatal(err)
		}

		elements = append(elements, e)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	elementsJSON, err := json.Marshal(elements)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, string(elementsJSON))
}

func add(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var element Element
	element.ID = primitive.NewObjectID()
	json.Unmarshal(reqBody, &element)

	_, err := collection.InsertOne(ctx, element)
	if err != nil {
		log.Fatal(err)
	}
}
