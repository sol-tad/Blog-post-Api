package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)
// SetupBlogRoutes registers routes related to blog posts.
// It initializes repositories for blogs, users, and interactions,
// creates the blog use case and controller, then sets up routes to
// list, get, create, update, and delete blog posts.
// Creation, update, and deletion routes are protected by authentication middleware.
func SetupBlogRoutes(router *gin.Engine) {
	// Initialize blog repository with blog collection
	blogRepo := repository.NewBlogRepo(config.BlogCollection)

	// Initialize user repository with user collection
	userRepo := repository.NewUserRepository(config.UserCollection)

	// Initialize interaction repository with blog and interaction collections
	interactionRepo := repository.NewInteractionRepository(
		config.BlogCollection,
		config.InteractionCollection,
	)

	// Create blog use case with blog, interaction, and user repositories
	blogUsecase := usecase.NewBlogUseCase(blogRepo, interactionRepo, userRepo)

	// Initialize blog controller with the use case
	blogController := controllers.NewBlogController(blogUsecase)

	// Group routes under /blogs
	blogRoutes := router.Group("/blogs")
	{
		// Public routes

		// List all blogs
		blogRoutes.GET("", blogController.ListBlogs)

		// Get a single blog by id
		blogRoutes.GET("/:id", blogController.GetBlog)

		// Protected routes requiring authentication
		protected := blogRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Create a new blog post
			protected.POST("/create", blogController.CreateBlog)

			// Update a blog post by id
			protected.PUT("/:id", blogController.UpdateBlog)

			// Delete a blog post by id
			protected.DELETE("/:id", blogController.DeleteBlog)
		}
	}
}
