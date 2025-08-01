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
}
