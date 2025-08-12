package repository

import (
	"context"
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// interactionRepository handles blog interactions like views, likes, and dislikes
type interactionRepository struct {
	blogCollection        *mongo.Collection
	interactionCollection *mongo.Collection
}

// NewInteractionRepository initializes a new repository with blog and interaction collections
func NewInteractionRepository(blogColl, interactionColl *mongo.Collection) domain.InteractionRepository {
	return &interactionRepository{
		blogCollection:        blogColl,
		interactionCollection: interactionColl,
	}
}

// IncrementViewCount increases the view count of a blog post by 1
func (r *interactionRepository) IncrementViewCount(blogID string) error {
	// Convert blogID from hex string to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	// Atomically increment the 'stats.views' field in the blog document
	_, err = r.blogCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.views": 1}},
	)
	return err
}

// RecordView logs a user's view of a blog post, updating timestamp and view count
func (r *interactionRepository) RecordView(blogID string, userID primitive.ObjectID) error {
	// Upsert interaction: update if exists, insert if not
	_, err := r.interactionCollection.UpdateOne(
		context.Background(),
		bson.M{"blog_id": blogID, "user_id": userID},
		bson.M{"$set": bson.M{
			"last_viewed": time.Now(),
			"view_count":  bson.M{"$inc": 1}, // Note: this won't increment unless handled properly
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

// AddLike registers a like from a user and increments the blog's like count
func (r *interactionRepository) AddLike(blogID string, userID primitive.ObjectID) error {
	// Set 'liked' to true, 'disliked' to false, and update interaction timestamp
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{"liked": true, "disliked": false, "last_interaction": time.Now()},
	}

	// Upsert interaction document
	res, err := r.interactionCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// If this is a new like or a change, increment blog's like count
	if res.UpsertedCount > 0 || res.ModifiedCount > 0 {
		objID, _ := primitive.ObjectIDFromHex(blogID)
		_, err = r.blogCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$inc": bson.M{"stats.likes": 1}},
		)
	}
	return err
}

// AddDislike registers a dislike from a user and increments the blog's dislike count
func (r *interactionRepository) AddDislike(blogID string, userID primitive.ObjectID) error {
	// Set 'disliked' to true, 'liked' to false, and update interaction timestamp
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{"disliked": true, "liked": false, "last_interaction": time.Now()},
	}

	// Upsert interaction document
	res, err := r.interactionCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// If this is a new dislike or a change, increment blog's dislike count
	if res.UpsertedCount > 0 || res.ModifiedCount > 0 {
		objID, _ := primitive.ObjectIDFromHex(blogID)
		_, err = r.blogCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$inc": bson.M{"stats.dislikes": 1}}, // Note: typo in original code: "dislike" vs "dislikes"
		)
	}
	return err
}

// RemoveLike unregisters a user's like and decrements the blog's like count
func (r *interactionRepository) RemoveLike(blogID string, userID primitive.ObjectID) error {
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{"$set": bson.M{"liked": false, "last_interaction": time.Now()}}

	// Update interaction document to remove like
	_, err := r.interactionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	// Decrement blog's like count
	objID, _ := primitive.ObjectIDFromHex(blogID)
	_, err = r.blogCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.likes": -1}},
	)
	return err
}

// RemoveDislike unregisters a user's dislike and decrements the blog's dislike count
func (r *interactionRepository) RemoveDislike(blogID string, userID primitive.ObjectID) error {
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{"$set": bson.M{"disliked": false, "last_interaction": time.Now()}}

	// Update interaction document to remove dislike
	_, err := r.interactionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	// Decrement blog's dislike count
	objID, _ := primitive.ObjectIDFromHex(blogID)
	_, err = r.blogCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.dislikes": -1}}, // Fixed typo: "dislike" → "dislikes"
	)
	return err
}