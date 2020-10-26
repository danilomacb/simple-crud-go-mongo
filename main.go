package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println("Disconnect")
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Error on Ping")
		panic(err)
	}

	collection = client.Database("simpleCrudGoMongo").Collection("elements")

	fmt.Println("Server running on port 3001")

	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", add).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/", list).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/{id}", delete).Methods("DELETE", "OPTIONS")
	myRouter.HandleFunc("/{id}", update).Methods("PUT", "OPTIONS")
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
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func list(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	var elements []Element

	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Println("Error on list, fail to find elements:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var e Element
		err := cur.Decode(&e)
		if err != nil {
			log.Println("Error on list, fail to decode element:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		elements = append(elements, e)
	}

	if err := cur.Err(); err != nil {
		log.Println("Error on list, something went wrong on findeds elements:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	elementsJSON, err := json.Marshal(elements)
	if err != nil {
		log.Println("Error on list, fail to convert to json:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintln(w, string(elementsJSON))
	w.WriteHeader(http.StatusOK)
}

func add(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error on add, fail to read request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var element Element
	element.ID = primitive.NewObjectID()

	err2 := json.Unmarshal(reqBody, &element)
	if err2 != nil {
		log.Println("Error on add, fail to convert json to Element:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err3 := collection.InsertOne(ctx, element)
	if err3 != nil {
		log.Println("Error on add, fail to insert one:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func delete(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error on delete, fail to transform id on primitive:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	filter := bson.M{"_id": idPrimitive}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error on delete, fail to delete one:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if res.DeletedCount == 0 {
		log.Println("Error on delete, no tasks were deleted")
	}

	w.WriteHeader(http.StatusOK)
}

func update(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	var e Element

	err2 := json.NewDecoder(r.Body).Decode(&e)
	if err2 != nil {
		log.Fatal(err2)
	}

	filter := bson.M{"_id": idPrimitive}
	update := bson.D{{"$set", bson.D{{"content", e.Content}}}}

	var updatedDocument bson.M
	err3 := collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)
	if err3 != nil {
		if err3 == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err3)
	}

	w.WriteHeader(http.StatusOK)
}
