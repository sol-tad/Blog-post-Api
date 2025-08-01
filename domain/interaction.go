package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// INTERACTION STRUCT

type UserInteraction struct {
	BlogID     primitive.ObjectID `bson:"blog_id"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Liked      bool               `bson:"liked,omitempty"`
	Disliked   bool               `bson:"disliked,omitempty"`
	LastViewed time.Time          `bson:"last_viewed,omitempty"`
}

// THIS IS INTERFACE FOR INTERACTION DATA ACTIONS 
type InteractionRepository interface {
	RecordView(blogID string, userID primitive.ObjectID) error
	AddLike(blogID string, userID primitive.ObjectID) error
	RemoveLike(blogID string, userID primitive.ObjectID) error
	AddDislike(blogID string, userID primitive.ObjectID) error
	RemoveDislike(blogID string, userID primitive.ObjectID) error
}
