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


func SetupAIRoutes(router *gin.Engine) {

	// userDbCollection:=config.UserCollection

	// userRepository:=repository.NewUserRepository(userDbCollection)
	// userUsecase:=usecase.NewUserUsecase(userRepository)
	// userController:=controllers.NewUserController(userUsecase)

	geminiApiKey := os.Getenv("GEMINI_API_KEY")
    log.Printf("------********************** %s  %s", geminiApiKey,"-------------------||")
	gemini,_:=infrastructure.NewGeminiService(geminiApiKey)
	aiUsecase:=usecase.NewAIUsecase(gemini)
	aiController:=controllers.NewAIController(aiUsecase)

	AIRoutes:=router.Group("/ai/posts")

	{
		
		AIRoutes.POST("generate", middlewares.AuthMiddleware(), aiController.GenerateBlogPost)
		AIRoutes.POST("/:id/improve", middlewares.AuthMiddleware(), aiController.ImproveBlogPost)
		AIRoutes.POST("/:id/suggestions", middlewares.AuthMiddleware(), aiController.SuggestBlogImprovements)
		AIRoutes.POST("/:id/summary", middlewares.AuthMiddleware(), aiController.SummarizeBlog)
		AIRoutes.POST("/:id/metadata", middlewares.AuthMiddleware(), aiController.GenerateMetadata)

	}
}



