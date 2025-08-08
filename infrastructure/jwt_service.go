package infrastructure

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secrets for signing JWT tokens, loaded from environment variables
var JWT_ACCESS_TOKEN_SECRET = os.Getenv("JWT_ACCESS_TOKEN_SECRET")
var JWT_REFRESH_TOKEN_SECRET = os.Getenv("JWT_REFRESH_TOKEN_SECRET")

// GenerateAccessToken creates a JWT access token with a short expiration (15 minutes)
// It includes user ID and role as claims
func GenerateAccessToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // expiration time
	}

	// Create a new token with HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and return it as a string
	return token.SignedString([]byte(JWT_ACCESS_TOKEN_SECRET))
}

// GenerateRefreshToken creates a JWT refresh token with a longer expiration (7 days)
// It includes only user ID as a claim
func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // expiration time
	}

	// Create a new token with HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the refresh token secret key and return it
	return token.SignedString([]byte(JWT_REFRESH_TOKEN_SECRET))
}

// VerifyRefreshToken validates a refresh token string and returns the user ID if valid
func VerifyRefreshToken(tokenStr string) (string, error) {
	// Parse the token with the refresh token secret key
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_REFRESH_TOKEN_SECRET), nil
	})

	// If there's an error or token is invalid, return empty string and error
	if err != nil || !token.Valid {
		return "", err
	}

	// Extract claims from the token
	claims := token.Claims.(jwt.MapClaims)

	// Return the user_id claim as a string
	return claims["user_id"].(string), nil
}
