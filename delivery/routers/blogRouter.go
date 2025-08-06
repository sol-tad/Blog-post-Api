package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

func SetupBlogRoutes(router *gin.Engine) {
	blogRepo := repository.NewBlogRepo(config.BlogCollection)
	interactionRepo := repository.NewInteractionRepository(
		config.BlogCollection, 
		config.InteractionCollection,
	)
	
	blogUsecase := usecase.NewBlogUseCase(blogRepo, interactionRepo)
	blogController := controllers.NewBlogController(blogUsecase)

	blogRoutes := router.Group("/blogs")
	{
		blogRoutes.GET("", blogController.ListBlogs)
		blogRoutes.GET("/:id", blogController.GetBlog)
		
		// Protected routes
		protected := blogRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/create", blogController.CreateBlog)
			protected.PUT("/:id", blogController.UpdateBlog)
			protected.DELETE("/:id", blogController.DeleteBlog)
		}
	}
}