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

type interactionRepository struct {
	blogCollection      *mongo.Collection
	interactionCollection *mongo.Collection
}

func NewInteractionRepository(blogColl, interactionColl *mongo.Collection) domain.InteractionRepository {
	return &interactionRepository{
		blogCollection:      blogColl,
		interactionCollection: interactionColl,
	}
}


func (r *interactionRepository) IncrementViewCount(blogID string) error {
	// Convert blogID to ObjectID
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	
	// Update blog view count
	_, err = r.blogCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.views": 1}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *interactionRepository) RecordView(blogID string, userID primitive.ObjectID) error {
	// Create or update user interaction
	_, err := r.interactionCollection.UpdateOne(
		context.Background(),
		bson.M{"blog_id": blogID, "user_id": userID},
		bson.M{"$set": bson.M{
			"last_viewed": time.Now(),
			"view_count":  bson.M{"$inc": 1},
		}},
		options.Update().SetUpsert(true),
	)
	return err
}





func (r *interactionRepository) AddLike(blogID string, userID primitive.ObjectID) error {
	// Check if user has already liked
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{"liked": true, "disliked": false, "last_interaction": time.Now()},
	}
	
	res, err := r.interactionCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	
	// If this is a new like (not an update)
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

func (r *interactionRepository) AddDislike(blogID string, userID primitive.ObjectID) error {
	// Check if user has already liked
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{"disliked": true, "liked": false, "last_interaction": time.Now()},
	}
	
	res, err := r.interactionCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	
	// If this is a new like (not an update)
	if res.UpsertedCount > 0 || res.ModifiedCount > 0 {
		objID, _ := primitive.ObjectIDFromHex(blogID)
		_, err = r.blogCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$inc": bson.M{"stats.dislikes": 1}},
		)
	}
	return err
}



func (r *interactionRepository) RemoveLike(blogID string, userID primitive.ObjectID) error {
	filter := bson.M{"blog_id": blogID, "user_id": userID}
	update := bson.M{"$set": bson.M{"liked": false, "last_interaction": time.Now()}}
	
	_, err := r.interactionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	
	objID, _ := primitive.ObjectIDFromHex(blogID)
	_, err = r.blogCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stats.likes": -1}},
	)
	return err
}

func (r *interactionRepository) RemoveDislike(blogID string, userID primitive.ObjectID) error {
	filter:= bson.M{"blog_id":blogID, "user_id":userID}
	update := bson.M{"$set":bson.M{"disliked":false ,"last_interaction":time.Now()}}
	_,err := r.interactionCollection.UpdateOne(context.Background(), filter, update) 
	if err != nil {
		return err 
	}
	objID,_ := primitive.ObjectIDFromHex(blogID)
	_,err = r.blogCollection.UpdateOne(
		context.Background(),
		bson.M{"_id":objID},
		bson.M{"$inc":bson.M{"stats.dislike":-1}},
	)
	return err 
}
