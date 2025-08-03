package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       	primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username 	string 		   `json:"username" bson:"username" validate:"required,min=3,max=50"`	
	Email      string             `bson:"email" json:"email"`
	Password 	string `json:"password" bson:"password" validate:"required,min=6,max=50"`
	Role 		string `json:"role" bson:"role"`
	RefreshToken string `bson:"refresh_token,omitempty"`
	OTPCode    string             `bson:"otp_code"`
	ResetOTP     string             `bson:"reset_otp"` 
	IsVerified bool               `bson:"is_verified"`
	Bio           string             `json:"bio,omitempty" bson:"bio,omitempty"`
	ProfilePicture string            `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
	ContactInfo   string             `json:"contact_info,omitempty" bson:"contact_info,omitempty"`

}

type UserRepository interface {
	Register(ctx context.Context,user User) (User, error)
	Login(ctx context.Context,username string) (User, error)
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


	// PromoteUser(ctx context.Context, adminID string, targetUserID string) error
	// DemoteUser(ctx context.Context, adminID string, targetUserID string) error

	UpdateUserRole(ctx context.Context, userID string, role string) error
}
