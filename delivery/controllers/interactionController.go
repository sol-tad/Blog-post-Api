package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InteractionController handles blog interaction endpoints (like, dislike, etc.)
type InteractionController struct {
	InteractionUsecase *usecase.InteractionUsecase
}

// NewInteractionController initializes a new InteractionController
func NewInteractionController(interactionUsecase *usecase.InteractionUsecase) *InteractionController {
	return &InteractionController{
		InteractionUsecase: interactionUsecase,
	}
}

// LikeBlog allows a user to like a blog post
func (ic *InteractionController) LikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID := c.GetString("id") // Extracted from middleware
	log.Println("id============:", userID)

	// Convert user ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Register like via usecase
	if err := ic.InteractionUsecase.InteractionRepo.AddLike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}

// UnlikeBlog removes a user's like from a blog post
func (ic *InteractionController) UnlikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Convert user ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Remove like via usecase
	if err := ic.InteractionUsecase.InteractionRepo.RemoveLike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}

// DislikeBlog allows a user to dislike a blog post
func (ic *InteractionController) DislikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Convert user ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Register dislike via usecase
	if err := ic.InteractionUsecase.InteractionRepo.AddDislike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog disliked successfully"})
}

// UndoDislike removes a user's dislike from a blog post
func (ic *InteractionController) UndoDislike(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Convert user ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Remove dislike via usecase
	if err := ic.InteractionUsecase.InteractionRepo.RemoveDislike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dislike removed successfully"})
}