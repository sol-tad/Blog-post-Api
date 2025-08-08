// Package config handles application configuration,
// including OAuth2 setup for Google authentication.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)
// GoogleOAuthConfig holds the OAuth2 configuration for Google sign-in.
var GoogleOAuthConfig *oauth2.Config

// init loads environment variables from a .env file (if available)
// and initializes GoogleOAuthConfig with client credentials, redirect URL,
// required scopes, and Google's OAuth2 endpoint.
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file")
	}

	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/oauth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
