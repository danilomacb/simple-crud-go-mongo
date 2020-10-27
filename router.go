package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func router() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", add).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/", list).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/{id}", delete).Methods("DELETE", "OPTIONS")
	myRouter.HandleFunc("/{id}", update).Methods("PUT", "OPTIONS")
	myRouter.Use(mux.CORSMethodMiddleware(myRouter))
	log.Fatal(http.ListenAndServe(":3001", myRouter))
}
