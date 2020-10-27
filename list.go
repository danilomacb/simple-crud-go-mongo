package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

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
}
