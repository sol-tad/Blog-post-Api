package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog represents a blog post entity stored in the database.
// swagger:model Blog
type Blog struct {
	// Unique identifier for the blog (MongoDB ObjectID)
	// example: 507f1f77bcf86cd799439011
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// Blog post title (required)
	// example: Introduction to Go
	Title string `json:"title,omitempty" bson:"title,omitempty" validate:"required"`

	// Blog content body (required)
	// example: This is an introductory post about the Go programming language...
	Content string `json:"content,omitempty" bson:"content,omitempty" validate:"required"`

	// ID of the author who wrote the blog (required)
	// example: 507f1f77bcf86cd799439022
	AuthorID primitive.ObjectID `json:"author_id" bson:"author_id" validate:"required"`

	// Author's name (required)
	// example: John Doe
	AuthorName string `json:"author_name" bson:"author_name" validate:"required"`

	// List of tags associated with the blog post (required)
	// example: ["programming", "golang", "tutorial"]
	Tags []string `json:"tags" bson:"tags" validate:"required"`

	// Timestamp when blog was created
	// example: 2025-08-10T15:04:05Z
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	// Timestamp when blog was last updated
	// example: 2025-08-11T16:05:06Z
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	// Engagement metrics (views, likes, etc.)
	Stats BlogStats `json:"stats" bson:"stats"`
}

// BlogStats holds engagement statistics related to a blog post.
// swagger:model BlogStats
type BlogStats struct {
	// Number of times the blog was viewed
	// example: 1234
	Views int `json:"views" bson:"views"`

	// Number of likes received
	// example: 256
	Likes int `json:"likes" bson:"likes"`

	// Number of dislikes received
	// example: 10
	Dislikes int `json:"dislikes" bson:"dislikes"`

	// Number of comments on the blog
	// example: 45
	Comments int `json:"comments" bson:"comments"`
}

// BlogRepositary defines the interface for data operations related to blogs.
// No swagger annotations needed for interfaces
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
// swagger:model BlogFilter
type BlogFilter struct {
	// Text search query to filter blogs by title/content
	// example: Go programming
	Search string

	// Filter blogs by author name or ID
	// example: John Doe
	Author string

	// Filter blogs by one or more tags
	// example: ["golang", "tutorial"]
	Tag []string

	// Filter blogs created after this date
	// example: 2025-01-01T00:00:00Z
	StartDate time.Time

	// Filter blogs created before this date
	// example: 2025-12-31T23:59:59Z
	EndDate time.Time

	// Field to sort by (e.g., "created_at", "views")
	// example: created_at
	SortBy string

	// Sort order, "asc" for ascending or "dsc" for descending
	// example: asc
	SortOrder string
}
