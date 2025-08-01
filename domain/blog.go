package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Blog struct {
	ID         primitive.ObjectID `bson:"_id"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	AuthorID   primitive.ObjectID `bson:"author_id"`
	AuthorName string             `bson:"author_name"`
	Tags       []string           `bson:"tags"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	Stats      BlogStats          `bson:"stats"`

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

