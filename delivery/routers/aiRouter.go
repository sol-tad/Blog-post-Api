package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

func SetupAI(router *gin.Engine) {
    // 1. Create the Gemini Adapter (infrastructure layer)
    adapter, err := infrastructure.NewGeminiAdapter()
    if err != nil {
        panic("Failed to initialize Gemini Adapter: " + err.Error())
    }

    // 2. Create the Use Case with the adapter (business logic layer)
    aiUC := usecase.NewAIUseCases(adapter)

    // 3. Create the Controller with the use case (interface layer)
    aiController := controllers.NewAIController(aiUC)

    // 4. Define route group for AI endpoints
    aiGroup := router.Group("/ai")
    {
        aiGroup.POST("/generate", aiController.GenerateBlog)
        aiGroup.POST("/summarize", aiController.SummarizeBlog)
    }
}
