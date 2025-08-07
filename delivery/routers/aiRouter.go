package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

func SetupAiRoutes(router *gin.Engine){
	AiUseCase := usecase.NewAiUseCase()
	AiController := controllers.NewAiController(AiUseCase)

	ai := router.Group("/ai")
	{
		ai.GET("/besttitle", AiController.BestTitleController)
		ai.GET("/emphasize", AiController.EmphasizeController)
		ai.GET("/summerize", AiController.SummerizeController)
		ai.GET("/improve", AiController.ImproveController)
	}	
}