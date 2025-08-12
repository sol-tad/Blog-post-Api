// Package routers initializes and configures all HTTP routes for the application.
package routers

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter creates a new Gin engine and registers all route groups,
// including user, AI, interaction (likes), OAuth, comment, and blog routes.
// It returns the configured *gin.Engine ready to run.
func SetupRouter() *gin.Engine{
	router:=gin.Default()
	
	 // Register user-related routes
    SetupUserRoutes(router)

	// Register AI-related routes
	SetupAIRoutes(router) 

	// Register interaction (like) related routes
	SetupInteractionRoutes(router)
	
	
	// Register OAuth authentication routes
	SetupOAuthRouter(router)

	// Register comment-related routes
	SetupCommentRoutes(router)

	// Register blog post related routes
	SetupBlogRoutes(router)
	
	return router
}