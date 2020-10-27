package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type element struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Content string             `bson:"content" json:"content"`
}
