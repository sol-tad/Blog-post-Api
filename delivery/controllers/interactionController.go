package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InteractionController struct {
	InteractionUsecase *usecase.InteractionUsecase
}

func NewInteractionController(interactionUsecase *usecase.InteractionUsecase) *InteractionController {
	return &InteractionController{
		InteractionUsecase: interactionUsecase,
	}
}

func (ic *InteractionController) LikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	// userID, exists := c.Get("id")

	userID := c.GetString("id")
	log.Println("id============:", userID)


	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
	// 	return
	// }

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := ic.InteractionUsecase.InteractionRepo.AddLike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}

func (ic *InteractionController) UnlikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := ic.InteractionUsecase.InteractionRepo.RemoveLike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}

func (ic *InteractionController) DislikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := ic.InteractionUsecase.InteractionRepo.AddDislike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog disliked successfully"})
}

func (ic *InteractionController) UndoDislike(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := ic.InteractionUsecase.InteractionRepo.RemoveDislike(blogID, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dislike removed successfully"})
}