package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupOAuthRouter configures OAuth-related routes on the provided Gin router.
// It initializes the necessary repository, use case, and controller for OAuth operations,
// then registers endpoints for Google OAuth login and callback handling.
func SetupOAuthRouter(router *gin.Engine) {
	// Retrieve the MongoDB user collection from config
	userDbCollection := config.UserCollection

	// Initialize user repository using the user collection
	userRepository := repository.NewUserRepository(userDbCollection)

	// Create the OAuth use case with the user repository
	oauthUsecase := usecase.NewOAuthUsecase(userRepository)

	// Initialize the OAuth controller with the use case
	oauthController := controllers.NewOAuthController(oauthUsecase)

	// Group routes without a specific prefix (root-level routes)
	oauthRoutes := router.Group("")

	{
		// Route to initiate Google OAuth login flow
		oauthRoutes.GET("/oauth/google/login", oauthController.Login)

		// OAuth callback endpoint to handle Google's response
		oauthRoutes.GET("/oauth/google/callback", oauthController.Callback)
	}
}
