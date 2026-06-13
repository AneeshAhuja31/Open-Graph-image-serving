package main

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Points int `json:"points" bson:"points"`
}

var Client *mongo.Client

func InitDB() {

	uri := os.Getenv("MONGODB_URI")
	if uri == ""{
		uri = "mongodb://localhost:27017"
	}
	ClientOptions := options.Client().ApplyURI(uri)
	var err error
	Client,err = mongo.Connect(ClientOptions)
	if err != nil{
		log.Fatalf("Failed connection with client")
	}	
}