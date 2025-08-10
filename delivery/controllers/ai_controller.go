package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

type AIController struct {
    aiUC *usecase.AIUseCase
}

func NewAIController(aiUC *usecase.AIUseCase) *AIController {
    return &AIController{aiUC: aiUC}
}

// Generate blog with parameters
func (c *AIController) GenerateBlog(ctx *gin.Context) {
    var params domain.GenerationParams
    
    if err := ctx.ShouldBindJSON(&params); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    content, err := c.aiUC.GenerateBlog(params)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed: " + err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, domain.AIResponse{Content: content})
}


// Summarize existing content
func (c *AIController) SummarizeBlog(ctx *gin.Context) {
    var request struct {
        Content string `json:"content" binding:"required"`
    }
    
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    summary, err := c.aiUC.SummarizeBLog(request.Content)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Summarization failed"})
        return
    }

    ctx.JSON(http.StatusOK, domain.AIResponse{Content: summary})
}