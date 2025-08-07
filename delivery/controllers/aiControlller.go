package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

type AiController struct{
		AiUseCase usecase.AiUseCase
}

func NewAiController (u usecase.AiUseCase) *AiController{
	return &AiController{
		AiUseCase: u,
	}
}

func (ac *AiController) SummerizeController (c *gin.Context){
	var req struct {
		Content string `json: "content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	summary , err := ac.AiUseCase.GenerateSummary(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK , gin.H{"summary" : summary})

}

func (ac *AiController)ImproveController (c *gin.Context){
	var req struct {
		Content string `json: "content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	summary , err := ac.AiUseCase.GenerateImproved(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK , gin.H{"Improved" : summary})

}

func (ac *AiController)EmphasizeController (c *gin.Context){
	var req struct {
		Content string `json: "content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	summary , err := ac.AiUseCase.GenerateEmphasized(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK , gin.H{"Emphasized" : summary})

}

func (ac *AiController)BestTitleController (c *gin.Context){
	var req struct {
		Content string `json: "content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	summary , err := ac.AiUseCase.GenerateTitle(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK , gin.H{"Best Title" : summary})

}

