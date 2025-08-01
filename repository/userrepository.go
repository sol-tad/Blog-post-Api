package repository

import (
	"context"
	"errors"
	"log"

	"github.com/sol-tad/Blog-post-Api/domian"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}
func NewUserRepository(coll *mongo.Collection) domian.UserRepository {
	return &UserRepositoryImpl{
		collection: coll,
	}
}


func (ur *UserRepositoryImpl) Register(ctx context.Context,user domian.User) (domian.User,error){

		user.Role="user"
		_,err:=ur.collection.InsertOne(ctx,user)
		return user,err
}
func (ur *UserRepositoryImpl) Login(ctx context.Context,username string)(domian.User,error){
	var user domian.User
    err := ur.collection.FindOne(context.Background(), map[string]string{"username": username}).Decode(&user)
        if err != nil {
        return user, errors.New("user not found")
    }
    return user, nil
}

func (ur *UserRepositoryImpl) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken}}

	_, err = ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to save refresh token")
	}
	return nil
}

func (ur *UserRepositoryImpl) VerifyRefreshToken(ctx context.Context, userID string, refreshToken string) (bool, error) {
	log.Println("idddd----------->:", userID)
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return false, errors.New("invalid user ID format")
    }

    var user domian.User
    err = ur.collection.FindOne(ctx, bson.M{
        "_id":           objID,
        "refresh_token":  refreshToken, // match the field name you used in SaveRefreshToken
    }).Decode(&user)

    if err != nil {
        return false, err
    }

    return true, nil
}



func (ur *UserRepositoryImpl) DeleteRefreshToken(ctx context.Context, userID string) error {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return errors.New("invalid user ID")
    }

    filter := bson.M{"_id": objID}
    update := bson.M{"$unset": bson.M{"refresh_token": ""}} 

    _, err = ur.collection.UpdateOne(ctx, filter, update)
    return err
}


func (ur *UserRepositoryImpl) FindByID(ctx context.Context, userID string) (domian.User, error) {
	var user domian.User

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user, errors.New("invalid user ID format")
	}

	err = ur.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}