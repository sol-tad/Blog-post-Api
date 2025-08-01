package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Blog struct {
	ID  primitive.ObjectID        `json:"id,omitempty" bson:"_id,omitempty"`
	Title string            	  `json:"title,omitempty" bson:"title,omitempty" validate:"required"`
	Content string 				  `json:"content,omitempty" bson:"content,omitempty" validate:"required"`
	AuthorID primitive.ObjectID   `json:"author_id" bson:"author_id" validate:"required"`
	AuthorName string             `json:"author_name" bson:"author_name" validate:"required"`
	Tags       []string           `json:"tags" bson:"tags" validate:"required"`
	CreatedAt   time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" bson:"updated_at"`
	Stats     BlogStats    		`json:"stats" bson:"stats"`



}

// BlogStats is the engagement metrics for blogs  
type BlogStats struct {
	Views int `json:"views" bson:"views"`
	Likes     int `json:"likes" bson:"likes"`
	Dislikes  int `json:"dislikes" bson:"dislikes"`
	Comments  int `json:"comments" bson:"comments"`
}



// This defines the interface for the blog data operations
type BlogRepositary interface {
	Create(blog *Blog) error
	GetByID( id string) (*Blog, error)
	GetByAuthor(author string, limt, page string) ([]*Blog, error)
	Update(blog *Blog) error
	Delete(id string) error
	List(page, limit int, filter BlogFilter) ([]*Blog, int64, error)
	IncrementViewCount(blogID string) error
	IncrementLikeCount(blogID string) error
	IncrementDislikeCount(blogID string) error
	DecrementLikeCount(blogID string) error
	DecrementDislikeCount(blogID string) error
	

}

// this contains parametres for filtering  
type BlogFilter struct {
	Search  	string
	 Author 	string
	 Tag    	string 
	 StratDate  time.Time 
	 EndDate    time.Time
	  SortBy    string
	  SortOrder string  // "asc" or "dsc"
}

