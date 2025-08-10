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
        // @Summary      Google OAuth Login
        // @Description  Initiates the Google OAuth 2.0 login flow.
        // @Tags         oauth
        // @Produce      json
        // @Success      302  {string}  string  "Redirect to Google's OAuth consent screen"
        // @Router       /oauth/google/login [get]
        oauthRoutes.GET("/oauth/google/login", oauthController.Login)

        // @Summary      Google OAuth Callback
        // @Description  Handles Google's OAuth 2.0 callback after user authorization.
        // @Tags         oauth
        // @Produce      json
        // @Param        state  query  string  true  "OAuth state parameter"
        // @Param        code   query  string  true  "Authorization code returned by Google"
        // @Success      200    {string}  string  "Authentication successful"
        // @Failure      400    {object}  map[string]string "Invalid or missing parameters"
        // @Router       /oauth/google/callback [get]
        oauthRoutes.GET("/oauth/google/callback", oauthController.Callback)
    }
}
