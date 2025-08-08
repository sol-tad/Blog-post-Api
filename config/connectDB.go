// Package config provides MongoDB connection setup and exposes
// global variables for accessing specific collections used in the app.

package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
// UserCollection represents the MongoDB collection for user documents.
var UserCollection *mongo.Collection

// UserCollection represents the MongoDB collection for blog post  documents.
var BlogCollection *mongo.Collection

// UserCollection represents the MongoDB collection for user interactions.
var InteractionCollection *mongo.Collection

// UserCollection represents the MongoDB collection for blog comments.
var CommentCollection *mongo.Collection


// ConnectDB connects to the MongoDB database using the connection string
// stored in the MONGODB_URI environment variable. It initializes the global
// collection variables for users, blogs, interactions, and comments within
// the "blogDB" database. If any error occurs during connection, the function
// logs the error and terminates the application.
func ConnectDB(){
    MONGODB_URI := os.Getenv("MONGODB_URI")
	client,err:=mongo.NewClient(options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err=client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	UserCollection=client.Database("blogDB").Collection("users")
	BlogCollection=client.Database("blogDB").Collection("blogs")
	InteractionCollection = client.Database("blogDB").Collection("interactions")
	CommentCollection = client.Database("blogDB").Collection("comments")
	log.Println("Connected to MongoDB")

}