package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the system.
type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName       string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Picture        string             `json:"picture,omitempty" bson:"picture,omitempty"`
	Username       string             `json:"username" bson:"username" validate:"required,min=3,max=50"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"password" bson:"password" validate:"required,min=6,max=50"`
	Role           string             `json:"role" bson:"role"`
	RefreshToken   string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	OTPCode        string             `json:"otp_code,omitempty" bson:"otp_code,omitempty"`
	ResetOTP       string             `json:"reset_otp,omitempty" bson:"reset_otp,omitempty"`
	IsVerified     bool               `json:"is_verified" bson:"is_verified"`
	Bio            string             `json:"bio,omitempty" bson:"bio,omitempty"`
	ProfilePicture string             `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
	ContactInfo    string             `json:"contact_info,omitempty" bson:"contact_info,omitempty"`
	GoogleID       string             `json:"google_id,omitempty" bson:"google_id,omitempty"`
}

// UserRepository defines methods for interacting with user data storage.
type UserRepository interface {
	Register(ctx context.Context, user User) (User, error)
	Login(ctx context.Context, username string) (User, error)
	SaveRefreshToken(ctx context.Context, userID string, token string) error
	VerifyRefreshToken(ctx context.Context, userID string, refreshToken string) (bool, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
	FindByID(ctx context.Context, userID string) (User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	VerifyUserOTP(ctx context.Context, email, otp string) error
	UpdateResetOTP(ctx context.Context, email, otp string) error
	VerifyResetOTP(ctx context.Context, email, otp string) error
	UpdatePasswordByEmail(ctx context.Context, email, newHashedPassword string) error
	UpdateProfile(ctx context.Context, userID string, updated User) (User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*User, error)
	Save(ctx context.Context, user *User) error
	GetByID(userID primitive.ObjectID) *User
	UpdateUserRole(ctx context.Context, userID string, role string) error
}
