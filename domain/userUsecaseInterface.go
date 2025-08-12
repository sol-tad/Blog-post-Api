package domain

import "context"

type UserUsecaseInterface interface {
	Register(ctx context.Context, user User) error
	VerifyOTP(ctx context.Context, email, otp string) error
	Login(ctx context.Context, username, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, userID string) error
	SendResetOTP(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, otp, newPassword string) error
	PromoteUser(ctx context.Context, adminID string, targetUserID string) error
	DemoteUser(ctx context.Context, adminID string, targetUserID string) error
	UpdateProfile(ctx context.Context, userID string, updated User) (User, error)
}