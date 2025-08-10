package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment on a blog post.
// swagger:model Comment
type Comment struct {
	// Unique identifier for the comment
	// example: 507f1f77bcf86cd799439011
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// ID of the blog post this comment belongs to
	// example: 507f1f77bcf86cd799439022
	BlogID primitive.ObjectID `json:"blog_id,omitempty" bson:"blog_id,omitempty"`

	// ID of the user who wrote the comment
	// example: 507f1f77bcf86cd799439033
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`

	// Username of the commenter
	// example: johndoe
	Username string `json:"username" bson:"username"`

	// Content/text of the comment
	// example: This is a great post! Thanks for sharing.
	Content string `json:"content" bson:"content"`

	// Timestamp when the comment was created
	// example: 2023-08-10T15:04:05Z07:00
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	// Timestamp when the comment was last updated
	// example: 2023-08-11T16:05:06Z07:00
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CommentRepository defines the interface for data operations on comments.
// No swagger annotations needed for interfaces
type CommentRepository interface {
	Create(comment *Comment) error                         // Create a new comment
	GetByID(id string) (*Comment, error)                   // Retrieve a comment by its ID
	GetByBlog(blogID string, page, limit int) ([]*Comment, error) // Retrieve comments for a blog with pagination
	Update(comment *Comment) error                          // Update an existing comment
	Delete(id string) error                                 // Delete a comment by its ID
	IncrementCommentCount(id string) error                  // Increment comment count for a related entity (e.g., blog)
	DecrementCommentCount(id string) error                  // Decrement comment count for a related entity
}
