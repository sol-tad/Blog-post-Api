package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domain"
)

// AIController handles AI-powered blog operations
type AIController struct {
	AIUsecase domain.AIUsecase
}

// NewAIController initializes a new AIController
func NewAIController(aiUsecase domain.AIUsecase) *AIController {
	return &AIController{AIUsecase: aiUsecase}
}

// GenerateBlogPost creates a blog post based on user input using AI
func (aicont *AIController) GenerateBlogPost(c *gin.Context) {
	var req domain.GenerateBlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := aicont.AIUsecase.GenerateBlogPost(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"generated_content": result})
}

// ImproveBlogPost enhances an existing blog post based on a specific goal
func (aicont *AIController) ImproveBlogPost(c *gin.Context) {
	var req domain.ImproveBlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := aicont.AIUsecase.ImproveBlogPost(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"improved_content": result})
}

// SuggestBlogImprovements provides writing and SEO suggestions for a blog post
func (aicont *AIController) SuggestBlogImprovements(c *gin.Context) {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suggestions, err := aicont.AIUsecase.SuggestBlogImprovement(body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

// SummarizeBlog generates a short summary of the blog post
func (aicont *AIController) SummarizeBlog(c *gin.Context) {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := aicont.AIUsecase.SummarizeBlog(body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

// GenerateMetadata extracts SEO metadata like title, tags, and description
func (aicont *AIController) GenerateMetadata(c *gin.Context) {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meta, err := aicont.AIUsecase.GenerateMetadata(body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metadata": meta})
}