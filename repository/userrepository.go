package repository

import (
	"context"
	"errors"
	"log"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}
func NewUserRepository(coll *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{
		collection: coll,
	}
}


func (ur *UserRepositoryImpl) Register(ctx context.Context, user domain.User) (domain.User, error) {
	_, err := ur.collection.InsertOne(ctx, user)
	return user, err
}

func (ur *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) VerifyUserOTP(ctx context.Context, email, otp string) error {
	filter := bson.M{"email": email, "otp_code": otp}
	update := bson.M{"$set": bson.M{"is_verified": true, "otp_code": ""}}

	res, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil || res.ModifiedCount == 0 {
		return errors.New("invalid email or OTP")
	}
	return nil
}

func (ur *UserRepositoryImpl) Login(ctx context.Context,username string)(domain.User,error){
	var user domain.User
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

    var user domain.User
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


func (ur *UserRepositoryImpl) FindByID(ctx context.Context, userID string) (domain.User, error) {
	var user domain.User

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

func (r *UserRepositoryImpl) UpdateResetOTP(ctx context.Context, email, otp string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"reset_otp": otp}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepositoryImpl) VerifyResetOTP(ctx context.Context, email, otp string) error {
	filter := bson.M{"email": email, "reset_otp": otp}

	var user domain.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return errors.New("invalid OTP or email")
	}
	return nil
}

func (r *UserRepositoryImpl) UpdatePasswordByEmail(ctx context.Context, email, newHashedPassword string) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"password":   newHashedPassword,
			"reset_otp":  "", // Clear OTP after success
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
