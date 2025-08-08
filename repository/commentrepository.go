package repository

import (
	"context"
	"errors"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// commentRepository implements the domain.CommentRepository interface
// and provides MongoDB-based operations for comments.
type commentRepository struct {
	collection *mongo.Collection // MongoDB collection for comments
}

// NewCommentRepository creates a new commentRepository instance
// using the provided MongoDB collection.
func NewCommentRepository(coll *mongo.Collection) domain.CommentRepository {
	return &commentRepository{
		collection: coll,
	}
}

// Create inserts a new comment document into the comments collection.
func (r *commentRepository) Create(comment *domain.Comment) error {
	_, err := r.collection.InsertOne(context.Background(), comment)
	return err
}

// GetByID finds and returns a comment by its string ID.
// Converts string ID to ObjectID internally.
func (r *commentRepository) GetByID(id string) (*domain.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var comment domain.Comment
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&comment)
	return &comment, err
}

// GetByBlog retrieves comments for a given blog ID, paginated by page and limit.
// Comments are sorted by creation time descending (most recent first).
func (r *commentRepository) GetByBlog(blogID string, page, limit int) ([]*domain.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}
	
	// Setup pagination and sorting options
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // sort descending by created_at
		SetSkip(int64((page - 1) * limit)).              // skip to the appropriate page
		SetLimit(int64(limit))                            // limit results to 'limit'

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"blog_id": objID},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode all matching comments into a slice
	var comments []*domain.Comment
	if err = cursor.All(context.Background(), &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// Update replaces an existing comment document identified by its ID with the provided comment data.
func (r *commentRepository) Update(comment *domain.Comment) error {
	_, err := r.collection.ReplaceOne(
		context.Background(),
		bson.M{"_id": comment.ID},
		comment,
	)
	return err
}

// Delete removes a comment document by its string ID.
func (r *commentRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}

// IncrementCommentCount increases the comment count field inside the stats sub-document
// for a blog identified by blogID. This function assumes the comment count is stored on the blog document.
func (r *commentRepository) IncrementCommentCount(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID")
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.comments": 1}}, // increment comments count by 1
	)
	return err
}

// DecrementCommentCount decreases the comment count for a blog by 1.
func (r *commentRepository) DecrementCommentCount(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID")
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.comments": -1}}, // decrement comments count by 1
	)
	return err
}
