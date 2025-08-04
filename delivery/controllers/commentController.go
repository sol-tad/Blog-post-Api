package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController struct {
	CommentUsecase *usecase.CommentUsecase
}

func NewCommentController(commentUsecase *usecase.CommentUsecase) *CommentController {
	return &CommentController{
		CommentUsecase: commentUsecase,
	}
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	var comment domain.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set user information
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	
	comment.UserID = objID
	comment.Username = c.GetString("username")
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	// Set blog ID
	blogID := c.Param("blog_id")
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}
	comment.BlogID = blogObjID

	if err := cc.CommentUsecase.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (cc *CommentController) GetComments(c *gin.Context) {
	blogID := c.Param("blog_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	comments, err := cc.CommentUsecase.GetCommentsByBlog(blogID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,

		},
	})
}

func (cc *CommentController) UpdateComment(c *gin.Context) {
	id := c.Param("id")
	
	var updatedComment domain.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get existing comment
	comment, err := cc.CommentUsecase.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}
	
	// Verify ownership
	userID := c.GetString("user_id")
	if comment.UserID.Hex() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own comments"})
		return
	}
	
	// Update content
	comment.Content = updatedComment.Content
	comment.UpdatedAt = time.Now()
	
	if err := cc.CommentUsecase.UpdateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, comment)
}

func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	
	// Get existing comment
	comment, err := cc.CommentUsecase.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}
	
	// Verify ownership or admin role
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")
	
	if comment.UserID.Hex() != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to delete this comment"})
		return
	}
	
	if err := cc.CommentUsecase.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}