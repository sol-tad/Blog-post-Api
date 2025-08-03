package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)

type UserUsecase struct {
	UserRepository domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase{
	return &UserUsecase{
		UserRepository: userRepo,
	}
}
func (uuc *UserUsecase) Register(ctx context.Context, user domain.User) error {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return errors.New("missing required fields")
	}

	// Check for existing user
	existing, _ := uuc.UserRepository.FindByEmail(ctx, user.Email)
	if existing != nil {
		return errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, _ := infrastructure.HashPassword(user.Password)
	user.Password = hashedPassword

	// Generate OTP
	user.OTPCode = fmt.Sprintf("%06d", rand.Intn(1000000))
	user.IsVerified = false
	user.Role = "user"

	// Save user
	_, err := uuc.UserRepository.Register(ctx, user)
	if err != nil {
		return err
	}

	// Send OTP via SendGrid
	return infrastructure.SendOTP(user.Email, user.OTPCode)
}

func (uuc *UserUsecase) VerifyOTP(ctx context.Context, email, otp string) error {
	return uuc.UserRepository.VerifyUserOTP(ctx, email, otp)
}



func (uuc *UserUsecase) Login(ctx context.Context, username, password string)(accessToken string, refreshToken string, err error) {
	user, err := uuc.UserRepository.Login(ctx, username)
	log.Println("USER&&&&&&&&&&&&&&&",user)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	//  Check if user is verified
	if !user.IsVerified {
		return "", "", errors.New("please verify your email before logging in")
	}

	if !infrastructure.CheckPassword(password, user.Password) {
		return "", "", errors.New("invalid username or password")
	}

	accessToken, err = infrastructure.GenerateAccessToken(user.ID.Hex(), user.Role)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = infrastructure.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return "", "", err
	}

	// Store the refresh token
	err = uuc.UserRepository.SaveRefreshToken(ctx, user.ID.Hex(), refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


func (uuc *UserUsecase) RefreshToken(ctx context.Context, refreshToken string)(string, error){


    userID, err := infrastructure.VerifyRefreshToken(refreshToken)
    if err != nil {
        return "", err
    }

	ok, err := uuc.UserRepository.VerifyRefreshToken(ctx, userID, refreshToken)
    if err != nil || !ok {
        return "", errors.New("invalid or expired refresh token")
    }

	    // Generate new access token
    user, err := uuc.UserRepository.FindByID(ctx, userID)
    if err != nil {
        return "", err
    }

    return infrastructure.GenerateAccessToken(user.ID.Hex(), user.Role)

}

func (uuc *UserUsecase) Logout(ctx context.Context, userID string) error {
    return uuc.UserRepository.DeleteRefreshToken(ctx, userID)
}


func (u *UserUsecase) SendResetOTP(ctx context.Context, email string) error {
	user, err := u.UserRepository.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	otp := infrastructure.GenerateOTP()
	if err := u.UserRepository.UpdateResetOTP(ctx, email, otp); err != nil {
		return err
	}

	return infrastructure.SendOTP(email, otp)
}

func (u *UserUsecase) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	err := u.UserRepository.VerifyResetOTP(ctx, email, otp)
	if err != nil {
		return err
	}

	hashedPassword, err := infrastructure.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return u.UserRepository.UpdatePasswordByEmail(ctx, email, hashedPassword)
}


func (uuc *UserUsecase) PromoteUser(ctx context.Context, adminID string, targetUserID string) error {
	admin, err := uuc.UserRepository.FindByID(ctx, adminID)
	if err != nil || admin.Role != "admin" {
		return errors.New("unauthorized")
	}
	return uuc.UserRepository.UpdateUserRole(ctx, targetUserID, "admin")
}

func (uuc *UserUsecase) DemoteUser(ctx context.Context, adminID string, targetUserID string) error {
	admin, err := uuc.UserRepository.FindByID(ctx, adminID)
	if err != nil || admin.Role != "admin" {
		return errors.New("unauthorized")
	}
	return uuc.UserRepository.UpdateUserRole(ctx, targetUserID, "user")
}

func (uuc *UserUsecase) UpdateProfile(ctx context.Context, userID string, updated domain.User) (domain.User, error) {
	return uuc.UserRepository.UpdateProfile(ctx, userID, updated)
}
