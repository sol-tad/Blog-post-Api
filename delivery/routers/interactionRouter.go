package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middleware"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

func SetupInteractionRoutes(router *gin.Engine) {
	interactionRepo := repository.NewInteractionRepository(
		config.BlogCollection, 
		config.InteractionCollection,
	)
	
	interactionUsecase := usecase.NewInteractionUsecase(interactionRepo)
	interactionController := controllers.NewInteractionController(interactionUsecase)

	interactionRoutes := router.Group("/blogs/:id")
	interactionRoutes.Use(middleware.AuthMiddleware())
	{
		interactionRoutes.PUT("/like", interactionController.LikeBlog)
		interactionRoutes.PUT("/unlike", interactionController.UnlikeBlog)
		interactionRoutes.PUT("/dislike", interactionController.DislikeBlog)
		interactionRoutes.PUT("/undislike", interactionController.UndoDislike)
	}
}