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

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(coll *mongo.Collection) domain.CommentRepository {
	return &commentRepository{
		collection: coll,
	}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	_, err := r.collection.InsertOne(context.Background(), comment)
	return err
}

func (r *commentRepository) GetByID(id string) (*domain.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var comment domain.Comment
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&comment)
	return &comment, err
}

func (r *commentRepository) GetByBlog(blogID string, page, limit int) ([]*domain.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page-1)*limit)).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"blog_id": objID},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	
	var comments []*domain.Comment
	if err = cursor.All(context.Background(), &comments); err != nil {
		return nil, err
	}
	
	return comments, nil
}

func (r *commentRepository) Update(comment *domain.Comment) error {
	_, err := r.collection.ReplaceOne(
		context.Background(),
		bson.M{"_id": comment.ID},
		comment,
	)
	return err
}

func (r *commentRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}



// IncrementCommentCount increases the comment count for a blog
func (r *commentRepository) IncrementCommentCount(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID")
	}
	
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.comments": 1}},
	)
	return err
}

// DecrementCommentCount decreases the comment count for a blog
func (r *commentRepository) DecrementCommentCount(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID")
	}
	
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.comments": -1}},
	)
	return err
}

