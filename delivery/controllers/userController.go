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

// VerifyOTPRequest defines payload for OTP verification
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// LoginRequest defines payload for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest defines payload for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// NewUserController initializes a new UserController
func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{UserUsecase: userUsecase}
}

// Register godoc
// @Summary Register a new user and send OTP for email verification
// @Description Registers a new user and sends an OTP to verify email
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User registration data"
// @Success 200 {object} map[string]string "OTP sent to your email. Please verify."
// @Failure 400 {object} map[string]string "Bad request error"
// @Router /users/register [post]
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

// VerifyOTP godoc
// @Summary Verify email using OTP
// @Description Verifies a user's email by OTP
// @Tags users
// @Accept json
// @Produce json
// @Param input body VerifyOTPRequest true "Email and OTP"
// @Success 200 {object} map[string]string "Account verified successfully!"
// @Failure 400 {object} map[string]string "Verification failed or bad input"
// @Router /users/verify-otp [post]
func (uc *UserController) VerifyOTP(c *gin.Context) {
	var input VerifyOTPRequest

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

// Login godoc
// @Summary Login user and return JWT tokens
// @Description Authenticates user credentials and returns access and refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "Access and refresh tokens"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /users/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest

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

// RefreshTokenController godoc
// @Summary Refresh access token
// @Description Generates a new access token using a valid refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string "New access token"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized - invalid token"
// @Router /users/refresh-token [post]
func (uc UserController) RefreshTokenController(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.UserUsecase.RefreshToken(context.Background(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidates user's refresh token to logout
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Logged out successfully"
// @Failure 500 {object} map[string]string "Logout failed"
// @Router /users/logout [post]
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

// SendResetOTP godoc
// @Summary Send password reset OTP
// @Description Sends a one-time password to user's email for password reset
// @Tags users
// @Accept json
// @Produce json
// @Param email body ForgotPasswordRequest true "Email to send OTP"
// @Success 200 {object} map[string]string "Reset OTP sent to your email address"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Failed to send reset OTP"
// @Router /users/send-reset-otp [post]
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

// ResetPassword godoc
// @Summary Reset user password
// @Description Verifies OTP and updates user's password
// @Tags users
// @Accept json
// @Produce json
// @Param resetRequest body ResetPasswordRequest true "Reset password details"
// @Success 200 {object} map[string]string "Password reset successfully"
// @Failure 400 {object} map[string]string "Invalid input or reset failed"
// @Router /users/reset-password [post]
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

// PromoteUser godoc
// @Summary Promote user to admin
// @Description Elevates a user's role to admin, requires admin privileges
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Target user ID"
// @Success 200 {object} map[string]string "User promoted to admin"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /users/promote/{id} [post]
func (uc *UserController) PromoteUser(c *gin.Context) {
	adminID := c.GetString("id") // from middleware
	targetID := c.Param("id")

	if err := uc.UserUsecase.PromoteUser(c, adminID, targetID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

// DemoteUser godoc
// @Summary Demote admin to regular user
// @Description Lowers a user's role to regular user, requires admin privileges
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Target user ID"
// @Success 200 {object} map[string]string "User demoted to regular user"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /users/demote/{id} [post]
func (uc *UserController) DemoteUser(c *gin.Context) {
	adminID := c.GetString("id") // from middleware
	targetID := c.Param("id")

	if err := uc.UserUsecase.DemoteUser(c, adminID, targetID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User demoted to regular user"})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Modifies the user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body domain.User true "Updated user profile data"
// @Success 200 {object} domain.User "Updated user profile"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Failed to update profile"
// @Router /users/profile [put]
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
