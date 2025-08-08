package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment on a blog post.
type Comment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`            // Unique identifier for the comment
	BlogID    primitive.ObjectID `json:"blog_id,omitempty" bson:"blog_id,omitempty"`    // ID of the blog post this comment belongs to
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`    // ID of the user who wrote the comment
	Username  string             `json:"username" bson:"username"`                        // Username of the commenter
	Content   string             `json:"content" bson:"content"`                          // Content/text of the comment
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`                    // Timestamp when the comment was created
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`                    // Timestamp when the comment was last updated
}

// CommentRepository defines the interface for data operations on comments.
type CommentRepository interface {
	Create(comment *Comment) error                         // Create a new comment
	GetByID(id string) (*Comment, error)                   // Retrieve a comment by its ID
	GetByBlog(blogID string, page, limit int) ([]*Comment, error) // Retrieve comments for a blog with pagination
	Update(comment *Comment) error                          // Update an existing comment
	Delete(id string) error                                 // Delete a comment by its ID
	IncrementCommentCount(id string) error                  // Increment comment count for a related entity (e.g., blog)
	DecrementCommentCount(id string) error                  // Decrement comment count for a related entity
}
