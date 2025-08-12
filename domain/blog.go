package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog represents a blog post entity stored in the database.
type Blog struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`             // Unique identifier for the blog (MongoDB ObjectID)
	Title      string             `json:"title,omitempty" bson:"title,omitempty" validate:"required"` // Blog post title (required)
	Content    string             `json:"content,omitempty" bson:"content,omitempty" validate:"required"` // Blog content body (required)
	AuthorID   primitive.ObjectID `json:"author_id" bson:"author_id" validate:"required"` // ID of the author who wrote the blog (required)
	AuthorName string             `json:"author_name" bson:"author_name" validate:"required"` // Author's name (required)
	Tags       []string           `json:"tags" bson:"tags" validate:"required"`            // List of tags associated with the blog post (required)
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`                     // Timestamp when blog was created
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`                     // Timestamp when blog was last updated
	Stats      BlogStats          `json:"stats" bson:"stats"`                               // Engagement metrics (views, likes, etc.)
}

// BlogStats holds engagement statistics related to a blog post.
type BlogStats struct {
	Views    int `json:"views" bson:"views"`       // Number of times the blog was viewed
	Likes    int `json:"likes" bson:"likes"`       // Number of likes received
	Dislikes int `json:"dislikes" bson:"dislikes"` // Number of dislikes received
	Comments int `json:"comments" bson:"comments"` // Number of comments on the blog
}

// BlogRepositary defines the interface for data operations related to blogs.
type BlogRepositary interface {
	Create(blog *Blog) error                                               // Create a new blog post
	GetByID(id string) (*Blog, error)                                     // Retrieve a blog post by its ID
	GetByAuthor(author string, limt, page string) ([]*Blog, error)        // Retrieve blogs by author with pagination
	Update(blog *Blog) error                                               // Update an existing blog post
	Delete(id string) error                                                // Delete a blog post by its ID
	List(page, limit int, filter BlogFilter) ([]*Blog, int64, error)      // List blogs with pagination and filtering
	IncrementViewCount(blogID string) error                               // Increment the view count of a blog
	IncrementLikeCount(blogID string) error                               // Increment the like count of a blog
	IncrementDislikeCount(blogID string) error                            // Increment the dislike count of a blog
	DecrementLikeCount(blogID string) error                               // Decrement the like count of a blog
	DecrementDislikeCount(blogID string) error                            // Decrement the dislike count of a blog
	IncrementCommentCount(blogID string) error                            // Increment the comment count of a blog
	DecrementCommentCount(blogID string) error                            // Decrement the comment count of a blog
}

// BlogFilter holds parameters used to filter and sort blog queries.
type BlogFilter struct {
	Search    string    // Text search query to filter blogs by title/content
	Author    string    // Filter blogs by author name or ID
	Tag       []string  // Filter blogs by one or more tags
	StartDate time.Time // Filter blogs created after this date
	EndDate   time.Time // Filter blogs created before this date
	SortBy    string    // Field to sort by (e.g., "created_at", "views")
	SortOrder string    // Sort order, "asc" for ascending or "dsc" for descending
}
