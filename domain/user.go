package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the system.
// swagger:model User
type User struct {
	// The unique identifier of the user
	// example: 507f1f77bcf86cd799439011
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// Full name of the user
	// example: John Doe
	FullName string `json:"full_name,omitempty" bson:"full_name,omitempty"`

	// URL to the user's profile picture
	// example: https://example.com/images/johndoe.jpg
	Picture string `json:"picture,omitempty" bson:"picture,omitempty"`

	// Username of the user
	// required: true
	// min length: 3
	// max length: 50
	// example: johndoe123
	Username string `json:"username" bson:"username" validate:"required,min=3,max=50"`

	// Email address of the user
	// example: johndoe@example.com
	Email string `json:"email" bson:"email"`

	// Password (hashed) for authentication
	// required: true
	// min length: 6
	// max length: 50
	// example: $2a$12$eXampleHashedPassword...
	Password string `json:"password" bson:"password" validate:"required,min=6,max=50"`

	// Role of the user (e.g., user, admin)
	// example: user
	Role string `json:"role" bson:"role"`

	// Refresh token for session management
	RefreshToken string `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`

	// OTP code for verification
	OTPCode string `json:"otp_code,omitempty" bson:"otp_code,omitempty"`

	// Reset OTP code for password reset
	ResetOTP string `json:"reset_otp,omitempty" bson:"reset_otp,omitempty"`

	// Whether the user has verified their email
	// example: true
	IsVerified bool `json:"is_verified" bson:"is_verified"`

	// User biography or description
	// example: Developer and blogger
	Bio string `json:"bio,omitempty" bson:"bio,omitempty"`

	// Profile picture URL
	// example: https://example.com/images/profilepic.jpg
	ProfilePicture string `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`

	// Contact information (phone, etc.)
	// example: +1234567890
	ContactInfo string `json:"contact_info,omitempty" bson:"contact_info,omitempty"`

	// Google account ID for OAuth login
	// example: 1234567890abcdefg
	GoogleID string `json:"google_id,omitempty" bson:"google_id,omitempty"`
}

// UserRepository defines methods for interacting with user data storage.
// (No swagger annotations needed for interfaces)
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
