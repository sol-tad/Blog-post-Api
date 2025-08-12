package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserInteraction represents a user's interaction with a blog post,
// including likes, dislikes, and last view time.
type UserInteraction struct {
	BlogID     primitive.ObjectID `bson:"blog_id"`                 // ID of the blog post
	UserID     primitive.ObjectID `bson:"user_id"`                 // ID of the user
	Liked      bool               `bson:"liked,omitempty"`         // Whether the user liked the post
	Disliked   bool               `bson:"disliked,omitempty"`      // Whether the user disliked the post
	LastViewed time.Time          `bson:"last_viewed,omitempty"`   // Timestamp of the last time user viewed the post
}

// InteractionRepository defines data operations related to user interactions with blogs.
type InteractionRepository interface {
	RecordView(blogID string, userID primitive.ObjectID) error  // Record a view event for a blog by a user
	AddLike(blogID string, userID primitive.ObjectID) error     // Add a like from a user to a blog
	RemoveLike(blogID string, userID primitive.ObjectID) error  // Remove a like from a user for a blog
	AddDislike(blogID string, userID primitive.ObjectID) error  // Add a dislike from a user to a blog
	RemoveDislike(blogID string, userID primitive.ObjectID) error // Remove a dislike from a user for a blog
	IncrementViewCount(blogID string) error                      // Increment the total view count of a blog post
}
