package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// OAuthController handles OAuth login and callback logic
type OAuthController struct {
	OAuthUsecase *usecase.OAuthUsecase
}

// NewOAuthController initializes a new OAuthController
func NewOAuthController(oauthUsecase *usecase.OAuthUsecase) *OAuthController {
	return &OAuthController{
		OAuthUsecase: oauthUsecase,
	}
}

// Login redirects the user to Google's OAuth consent screen
func (oauc *OAuthController) Login(c *gin.Context) {
	// Generate OAuth URL with a state parameter
	url := config.GoogleOAuthConfig.AuthCodeURL("state_value")

	// Redirect user to Google's OAuth page
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the OAuth redirect from Google and completes login
func (oauc *OAuthController) Callback(c *gin.Context) {
	// Extract authorization code from query parameters
	code := c.Query("code")
	log.Println("**************----", code, "-----------------------****")

	// Exchange code for tokens and user info
	user, accessToken, refreshToken, err := oauc.OAuthUsecase.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OAuth failed"})
		return
	}

	// Respond with user info and tokens
	c.JSON(http.StatusOK, gin.H{
		"message":       "OAuth login successful",
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}