package routers

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	// "github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
	"github.com/sol-tad/Blog-post-Api/middlewares"

	// "github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupAIRoutes configures AI-related routes for blog posts,
// initializing the AI service with API key, use case, and controller,
// then setting up authenticated endpoints for AI-powered blog operations.
func SetupAIRoutes(router *gin.Engine) {
	// Load the Gemini API key from environment variables
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	log.Printf("------********************** %s  %s", geminiApiKey, "-------------------||")

	// Initialize Gemini AI service with the API key
	gemini, _ := infrastructure.NewGeminiService(geminiApiKey)

	// Create AI use case with the Gemini service
	aiUsecase := usecase.NewAIUsecase(gemini)

	// Initialize AI controller with the use case
	aiController := controllers.NewAIController(aiUsecase)

	// Group routes under /ai/posts for AI-related blog post endpoints
	AIRoutes := router.Group("/ai/posts")

	{
		// Protected route to generate a new blog post using AI
		AIRoutes.POST("generate", middlewares.AuthMiddleware(), aiController.GenerateBlogPost)

		// Protected route to improve an existing blog post by ID
		AIRoutes.POST("/:id/improve", middlewares.AuthMiddleware(), aiController.ImproveBlogPost)

		// Protected route to get AI suggestions for blog improvements by ID
		AIRoutes.POST("/:id/suggestions", middlewares.AuthMiddleware(), aiController.SuggestBlogImprovements)

		// Protected route to get a summary of a blog post by ID
		AIRoutes.POST("/:id/summary", middlewares.AuthMiddleware(), aiController.SummarizeBlog)

		// Protected route to generate metadata for a blog post by ID
		AIRoutes.POST("/:id/metadata", middlewares.AuthMiddleware(), aiController.GenerateMetadata)
	}
}
