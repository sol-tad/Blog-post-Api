package routers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupAIRoutes configures AI-related routes for blog posts,
// initializing the AI service with API key, use case, and controller,
// then setting up authenticated endpoints for AI-powered blog operations.
func SetupAIRoutes(router *gin.Engine) {
	geminiApiKey := os.Getenv("GEMINI_API_KEY")

	gemini, _ := infrastructure.NewGeminiService(geminiApiKey)
	aiUsecase := usecase.NewAIUsecase(gemini)
	aiController := controllers.NewAIController(aiUsecase)

	AIRoutes := router.Group("/ai/posts")
	{
		// @Summary      Generate AI-powered blog post
		// @Description  Generates a new blog post using AI based on provided prompt.
		// @Tags         AI
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request body      map[string]interface{} true "Prompt data"
		// @Success      201  {object}  map[string]interface{}
		// @Failure      400  {object}  map[string]string
		// @Failure      401  {object}  map[string]string
		// @Router       /ai/posts/generate [post]
		AIRoutes.POST("generate", middlewares.AuthMiddleware(), aiController.GenerateBlogPost)

		// @Summary      Improve blog post
		// @Description  Improves an existing blog post by ID using AI suggestions.
		// @Tags         AI
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id      path      string                true  "Blog ID"
		// @Param        request body      map[string]interface{} true  "Improvement data"
		// @Success      200  {object}  map[string]interface{}
		// @Failure      400  {object}  map[string]string
		// @Failure      401  {object}  map[string]string
		// @Failure      404  {object}  map[string]string
		// @Router       /ai/posts/{id}/improve [post]
		AIRoutes.POST("/:id/improve", middlewares.AuthMiddleware(), aiController.ImproveBlogPost)

		// @Summary      AI suggestions for blog
		// @Description  Provides AI-generated suggestions for improving a blog post by ID.
		// @Tags         AI
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id   path      string  true  "Blog ID"
		// @Success      200  {object}  map[string]interface{}
		// @Failure      401  {object}  map[string]string
		// @Failure      404  {object}  map[string]string
		// @Router       /ai/posts/{id}/suggestions [post]
		AIRoutes.POST("/:id/suggestions", middlewares.AuthMiddleware(), aiController.SuggestBlogImprovements)

		// @Summary      Summarize blog post
		// @Description  Generates an AI-powered summary for a blog post by ID.
		// @Tags         AI
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id   path      string  true  "Blog ID"
		// @Success      200  {object}  map[string]interface{}
		// @Failure      401  {object}  map[string]string
		// @Failure      404  {object}  map[string]string
		// @Router       /ai/posts/{id}/summary [post]
		AIRoutes.POST("/:id/summary", middlewares.AuthMiddleware(), aiController.SummarizeBlog)

		// @Summary      Generate metadata for blog
		// @Description  Creates AI-generated metadata for a blog post by ID.
		// @Tags         AI
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id   path      string  true  "Blog ID"
		// @Success      200  {object}  map[string]interface{}
		// @Failure      401  {object}  map[string]string
		// @Failure      404  {object}  map[string]string
		// @Router       /ai/posts/{id}/metadata [post]
		AIRoutes.POST("/:id/metadata", middlewares.AuthMiddleware(), aiController.GenerateMetadata)
	}
}
