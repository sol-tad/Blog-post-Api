package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

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
	log.Println("Connected to MongoDB")

}