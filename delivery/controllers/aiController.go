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

// GenerateBlogPost godoc
// @Summary Generate a blog post based on user input using AI
// @Description Create a blog post from a request with AI assistance
// @Tags AI
// @Accept json
// @Produce json
// @Param request body domain.GenerateBlogPostRequest true "Blog generation request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/generate [post]
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

// ImproveBlogPost godoc
// @Summary Improve an existing blog post based on a goal
// @Description Enhance blog content using AI with a specified goal
// @Tags AI
// @Accept json
// @Produce json
// @Param request body domain.ImproveBlogPostRequest true "Blog improvement request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/improve [post]
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

// SuggestBlogImprovements godoc
// @Summary Suggest writing and SEO improvements for blog content
// @Description Get AI-generated suggestions to improve blog content quality and SEO
// @Tags AI
// @Accept json
// @Produce json
// @Param content body object true "Blog content"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/suggest [post]
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

// SummarizeBlog godoc
// @Summary Generate a short summary of the blog content
// @Description Get a concise AI-generated summary of blog content
// @Tags AI
// @Accept json
// @Produce json
// @Param content body object true "Blog content"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/summarize [post]
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

// GenerateMetadata godoc
// @Summary Generate SEO metadata from blog content
// @Description Extract title, tags, and description metadata using AI
// @Tags AI
// @Accept json
// @Produce json
// @Param content body object true "Blog content"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/metadata [post]
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
