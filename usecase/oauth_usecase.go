package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)

// OAuthUsecase handles OAuth-based authentication and user management
type OAuthUsecase struct {
	userRepo domain.UserRepository
}

// NewOAuthUsecase initializes a new OAuthUsecase with a user repository
func NewOAuthUsecase(urepo domain.UserRepository) *OAuthUsecase {
	return &OAuthUsecase{userRepo: urepo}
}

// HandleGoogleCallback processes the OAuth callback from Google
func (u *OAuthUsecase) HandleGoogleCallback(code string) (*domain.User, string, string, error) {
	// Exchange authorization code for access token
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", "", err
	}

	// Create an authenticated HTTP client using the token
	client := config.GoogleOAuthConfig.Client(context.Background(), token)

	// Fetch user info from Google API
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", "", err
	}
	defer response.Body.Close()

	// Read and parse the response body
	body, _ := io.ReadAll(response.Body)
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, "", "", err
	}

	// Check if user already exists in the database
	existingUser, _ := u.userRepo.FindByGoogleID(context.Background(), googleUser.ID)
	if existingUser == nil {
		// Create a new user if not found
		existingUser = &domain.User{
			GoogleID: googleUser.ID,
			Email:    googleUser.Email,
			FullName: googleUser.Name,
			Picture:  googleUser.Picture,
		}
	}

	// Generate access and refresh tokens
	accessToken, err := infrastructure.GenerateAccessToken(existingUser.ID.Hex(), existingUser.Role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := infrastructure.GenerateRefreshToken(existingUser.ID.Hex())
	if err != nil {
		return nil, "", "", err
	}

	// Save refresh token to user and persist in database
	existingUser.RefreshToken = refreshToken
	if err := u.userRepo.Save(context.Background(), existingUser); err != nil {
		return nil, "", "", err
	}

	// Return user and tokens
	return existingUser, accessToken, refreshToken, nil
}