package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

func SetupCommentRoutes(router *gin.Engine) {
	commentRepo := repository.NewCommentRepository(config.CommentCollection)
	blogRepo := repository.NewBlogRepo(config.BlogCollection)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo)
	commentController := controllers.NewCommentController(commentUsecase)

	commentRoutes := router.Group("/blogs/comments/:blog_id/")
	{
		commentRoutes.GET("", commentController.GetComments)
		
		// Protected routes
		protected := commentRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("", commentController.CreateComment)
			protected.PUT("/:comment_id", commentController.UpdateComment)
			protected.DELETE("/:comment_id", commentController.DeleteComment)
		}
	}
}