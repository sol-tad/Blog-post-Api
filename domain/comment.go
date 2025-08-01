package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// COMMENT STRUCT

type Comment struct {
	ID     primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty"`
	BlogID primitive.ObjectID   	  `json:"blog_id,omitempty" bson:"blog_id,omitempty"`
	UserID primitive.ObjectID         `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Username string 				  `json:"username" bson:"username"`
	Content string 					  `json:"content" bson:"content"`
	CreatedAt time.Time         	  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         	  `json:"updated_at" bson:"updated_at"`               
}

// THIS IS THE INTERFACE FOR COMMENT DATA OPERATIONS 
type CommentRepository interface {
	Create(comment *Comment) error
	GetByID(id string) (*Comment, error)
	GetByBlog( blogID string, page, limit int) ( *[]Comment, error)
	Update(comment *Comment) error 
	Delete(id string) error 
}
