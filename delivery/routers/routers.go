// Package routers initializes and configures all HTTP routes for the application.
package routers

import (
    "github.com/gin-gonic/gin"
)

// SetupRouter godoc
// @title           Blog API
// @version         1.0
// @description     This is the Blog API documentation.
// @host      localhost:8080
// @BasePath  /api/v1
func SetupRouter() *gin.Engine {
    router := gin.Default()



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
