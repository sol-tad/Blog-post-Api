package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)
// SetupInteractionRoutes registers routes related to blog post interactions,
// such as liking and disliking posts. It initializes the interaction repository,
// use case, and controller, then sets up authenticated routes under the
// "/blogs/:id" path to handle user interactions on specific blog posts.
func SetupInteractionRoutes(router *gin.Engine) {
	// Initialize the interaction repository with blog and interaction collections
	interactionRepo := repository.NewInteractionRepository(
		config.BlogCollection,
		config.InteractionCollection,
	)

	// Create the interaction use case with the repository
	interactionUsecase := usecase.NewInteractionUsecase(interactionRepo)

	// Initialize the interaction controller with the use case
	interactionController := controllers.NewInteractionController(interactionUsecase)

	// Create a route group for routes under /blogs/:id, representing individual blog posts
	interactionRoutes := router.Group("/blogs/:id")

	// Protect all interaction routes with authentication middleware
	interactionRoutes.Use(middlewares.AuthMiddleware())

	{
		// Route to like a blog post
		interactionRoutes.PUT("/like", interactionController.LikeBlog)

		// Route to unlike a previously liked blog post
		interactionRoutes.PUT("/unlike", interactionController.UnlikeBlog)

		// Route to dislike a blog post
		interactionRoutes.PUT("/dislike", interactionController.DislikeBlog)

		// Route to remove a previously given dislike
		interactionRoutes.PUT("/undislike", interactionController.UndoDislike)
	}
}
