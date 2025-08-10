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

// CreateComment godoc
// @Summary Add a new comment to a blog post
// @Description Create a comment on a blog post. Requires authentication.
// @Tags Comments
// @Accept json
// @Produce json
// @Param blog_id path string true "Blog ID"
// @Param comment body domain.Comment true "Comment object"
// @Success 201 {object} domain.Comment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /blogs/{blog_id}/comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	var comment domain.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := cc.CommentUsecase.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments godoc
// @Summary Get paginated comments for a blog post
// @Description Retrieve comments with pagination for a specific blog post
// @Tags Comments
// @Produce json
// @Param blog_id path string true "Blog ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of comments per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /blogs/{blog_id}/comments [get]
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

// UpdateComment godoc
// @Summary Update an existing comment
// @Description Update a comment if the user is the owner
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Param comment body domain.Comment true "Updated comment content"
// @Success 200 {object} domain.Comment
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [put]
func (cc *CommentController) UpdateComment(c *gin.Context) {
	id := c.Param("comment_id")

	var updatedComment domain.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := cc.CommentUsecase.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	userID := c.GetString("id")
	if comment.UserID.Hex() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own comments"})
		return
	}

	comment.Content = updatedComment.Content
	comment.UpdatedAt = time.Now()

	if err := cc.CommentUsecase.UpdateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment if the user is the owner or admin
// @Tags Comments
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [delete]
func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("comment_id")

	comment, err := cc.CommentUsecase.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	userID := c.GetString("id")
	userRole := c.GetString("role")

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
