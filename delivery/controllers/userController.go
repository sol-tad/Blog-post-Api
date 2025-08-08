package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// UserController handles HTTP requests related to user operations
type UserController struct {
	UserUsecase *usecase.UserUsecase
}

// ForgotPasswordRequest defines the payload for requesting a password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest defines the payload for resetting a password
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// NewUserController initializes a new UserController
func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{UserUsecase: userUsecase}
}

// Register handles user registration and sends OTP for email verification
func (uc *UserController) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserUsecase.Register(context.Background(), user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email. Please verify."})
}

// VerifyOTP confirms user's email using the OTP
func (uc *UserController) VerifyOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserUsecase.VerifyOTP(context.Background(), input.Email, input.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully!"})
}

// Login authenticates the user and returns access and refresh tokens
func (uc *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := uc.UserUsecase.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshTokenController issues a new access token using a valid refresh token
func (uc UserController) RefreshTokenController(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.UserUsecase.RefreshToken(context.Background(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token*****"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// Logout invalidates the user's refresh token
func (uc UserController) Logout(c *gin.Context) {
	userID := c.GetString("id") // Extracted from middleware
	log.Println("id============:", userID)

	err := uc.UserUsecase.Logout(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// SendResetOTP sends a password reset OTP to the user's email
func (uc *UserController) SendResetOTP(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Forgot password request received for:", req.Email)

	err := uc.UserUsecase.SendResetOTP(c, req.Email)
	if err != nil {
		log.Println("SendResetOTP error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset OTP sent to your email address"})
}

// ResetPassword verifies OTP and updates the user's password
func (uc *UserController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.UserUsecase.ResetPassword(c, req.Email, req.OTP, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// PromoteUser elevates a user's role to admin (requires admin privileges)
func (uc *UserController) PromoteUser(c *gin.Context) {
	adminID := c.GetString("id") // from middleware
	targetID := c.Param("id")

	if err := uc.UserUsecase.PromoteUser(c, adminID, targetID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

// DemoteUser lowers a user's role to regular user (requires admin privileges)
func (uc *UserController) DemoteUser(c *gin.Context) {
	adminID := c.GetString("id") // from middleware
	targetID := c.Param("id")

	if err := uc.UserUsecase.DemoteUser(c, adminID, targetID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User demoted to regular user"})
}

// UpdateProfile modifies the user's profile information
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("id") // from AuthMiddleware

	var updated domain.User
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.UserUsecase.UpdateProfile(c, userID, updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}