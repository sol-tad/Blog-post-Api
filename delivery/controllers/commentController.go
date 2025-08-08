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

// CommentController handles HTTP requests related to blog comments
type CommentController struct {
	CommentUsecase *usecase.CommentUsecase
}

// NewCommentController initializes a new CommentController
func NewCommentController(commentUsecase *usecase.CommentUsecase) *CommentController {
	return &CommentController{CommentUsecase: commentUsecase}
}

// CreateComment adds a new comment to a blog post
func (cc *CommentController) CreateComment(c *gin.Context) {
	var comment domain.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from context
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Set user and blog metadata
	comment.UserID = objID
	comment.Username = c.GetString("username")
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	blogID := c.Param("blog_id")
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}
	comment.BlogID = blogObjID

	// Save comment
	if err := cc.CommentUsecase.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments retrieves paginated comments for a specific blog post
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

// UpdateComment modifies an existing comment if the user is the owner
func (cc *CommentController) UpdateComment(c *gin.Context) {
	id := c.Param("comment_id")

	var updatedComment domain.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve existing comment
	comment, err := cc.CommentUsecase.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	// Verify ownership
	userID := c.GetString("id")
	if comment.UserID.Hex() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own comments"})
		return
	}

	// Apply update
	comment.Content = updatedComment.Content
	comment.UpdatedAt = time.Now()

	if err := cc.CommentUsecase.UpdateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment removes a comment if the user is the owner or an admin
func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("comment_id")

	// Retrieve existing comment
	comment, err := cc.CommentUsecase.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	// Verify ownership or admin role
	userID := c.GetString("id")
	userRole := c.GetString("role")

	if comment.UserID.Hex() != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to delete this comment"})
		return
	}

	// Delete comment
	if err := cc.CommentUsecase.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}