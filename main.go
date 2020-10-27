package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	connect()
	router()
}
