package repository

import (
	"context"
	"errors"
	"log"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepositoryImpl provides MongoDB-based implementation of UserRepository
type UserRepositoryImpl struct {
	collection *mongo.Collection
}

// NewUserRepository initializes a new UserRepository with the given collection
func NewUserRepository(coll *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{
		collection: coll,
	}
}

// Register inserts a new user document into the collection
func (ur *UserRepositoryImpl) Register(ctx context.Context, user domain.User) (domain.User, error) {
	_, err := ur.collection.InsertOne(ctx, user)
	return user, err
}

// FindByEmail retrieves a user by their email address
func (ur *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyUserOTP checks if the provided OTP matches and marks the user as verified
func (ur *UserRepositoryImpl) VerifyUserOTP(ctx context.Context, email, otp string) error {
	filter := bson.M{"email": email, "otp_code": otp}
	update := bson.M{"$set": bson.M{"is_verified": true, "otp_code": ""}}

	res, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil || res.ModifiedCount == 0 {
		return errors.New("invalid email or OTP")
	}
	return nil
}

// Login retrieves a user by username
func (ur *UserRepositoryImpl) Login(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}

// SaveRefreshToken stores a refresh token for the user
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

// VerifyRefreshToken checks if the provided refresh token matches the stored one
func (ur *UserRepositoryImpl) VerifyRefreshToken(ctx context.Context, userID string, refreshToken string) (bool, error) {
	log.Println("idddd----------->:", userID) // Consider removing or refining this log for production

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user ID format")
	}

	var user domain.User
	err = ur.collection.FindOne(ctx, bson.M{
		"_id":           objID,
		"refresh_token": refreshToken,
	}).Decode(&user)

	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteRefreshToken removes the refresh token from the user's document
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

// FindByID retrieves a user by their ObjectID
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

// UpdateResetOTP sets a new OTP for password reset
func (ur *UserRepositoryImpl) UpdateResetOTP(ctx context.Context, email, otp string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"reset_otp": otp}}

	_, err := ur.collection.UpdateOne(ctx, filter, update)
	return err
}

// VerifyResetOTP checks if the provided reset OTP matches the stored one
func (ur *UserRepositoryImpl) VerifyResetOTP(ctx context.Context, email, otp string) error {
	filter := bson.M{"email": email, "reset_otp": otp}

	var user domain.User
	err := ur.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return errors.New("invalid OTP or email")
	}
	return nil
}

// UpdatePasswordByEmail updates the user's password and clears the reset OTP
func (ur *UserRepositoryImpl) UpdatePasswordByEmail(ctx context.Context, email, newHashedPassword string) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"password":  newHashedPassword,
			"reset_otp": "", // Clear OTP after successful password reset
		},
	}

	_, err := ur.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateUserRole changes the user's role (e.g., admin, editor, viewer)
func (ur *UserRepositoryImpl) UpdateUserRole(ctx context.Context, userID string, role string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"role": role}}

	_, err = ur.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateProfile modifies user's bio, profile picture, and contact info
func (ur *UserRepositoryImpl) UpdateProfile(ctx context.Context, userID string, updated domain.User) (domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, err
	}

	update := bson.M{
		"$set": bson.M{
			"bio":             updated.Bio,
			"profile_picture": updated.ProfilePicture,
			"contact_info":    updated.ContactInfo,
		},
	}

	_, err = ur.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return domain.User{}, err
	}

	// Return updated user
	return ur.FindByID(ctx, userID)
}

//
// OAuth-related repository methods
//

// FindByGoogleID retrieves a user by their Google ID
func (ur *UserRepositoryImpl) FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Save upserts a user document based on Google ID
func (ur *UserRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
	filter := bson.M{"google_id": user.GoogleID}
	update := bson.M{"$set": user}
	opts := options.Update().SetUpsert(true)

	_, err := ur.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetByID retrieves a user by ObjectID without context (consider adding context for consistency)
func (ur *UserRepositoryImpl) GetByID(userID primitive.ObjectID) *domain.User {
	var user domain.User
	filter := bson.M{"_id": userID}
	err := ur.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return &domain.User{}
	}
	return &user
}