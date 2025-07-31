package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID 		primitive.ObjectID
	Title   string
	Content string
	Tags    string
	Date    time.Time
	User 	User

}