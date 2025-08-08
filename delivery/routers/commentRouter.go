package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)
// SetupCommentRoutes registers routes related to blog comments.
// It sets up the comment and blog repositories, the comment use case, and controller,
// then creates routes to get, create, update, and delete comments on blog posts.
// Routes for creating, updating, and deleting comments are protected by authentication middleware.
func SetupCommentRoutes(router *gin.Engine) {
	// Initialize the comment repository with the comment collection
	commentRepo := repository.NewCommentRepository(config.CommentCollection)

	// Initialize the blog repository with the blog collection (used for validation, etc.)
	blogRepo := repository.NewBlogRepo(config.BlogCollection)

	// Create the comment use case with comment and blog repositories
	commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo)

	// Initialize the comment controller with the use case
	commentController := controllers.NewCommentController(commentUsecase)

	// Group routes under /blogs/comments/:blog_id to handle comments for a specific blog post
	commentRoutes := router.Group("/blogs/comments/:blog_id/")
	{
		// Public route: Get all comments for the specified blog post
		commentRoutes.GET("", commentController.GetComments)

		// Protected routes requiring authentication
		protected := commentRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Create a new comment on the blog post
			protected.POST("", commentController.CreateComment)

			// Update an existing comment by comment_id
			protected.PUT("/:comment_id", commentController.UpdateComment)

			// Delete a comment by comment_id
			protected.DELETE("/:comment_id", commentController.DeleteComment)
		}
	}
}
