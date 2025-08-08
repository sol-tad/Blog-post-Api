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

// BlogController handles HTTP requests related to blog operations
type BlogController struct {
	BlogUsecase *usecase.BlogUseCase
}

// NewBlogController initializes a new BlogController
func NewBlogController(blogUsecase *usecase.BlogUseCase) *BlogController {
	return &BlogController{BlogUsecase: blogUsecase}
}

// CreateBlog handles blog creation by authenticated users
func (bc *BlogController) CreateBlog(c *gin.Context) {
	var blog domain.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
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

	// Set author and timestamps
	blog.AuthorID = objID
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	// Store blog
	if err := bc.BlogUsecase.StoreBlog(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

// GetBlog retrieves a blog post by its ID
func (bc *BlogController) GetBlog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing blog ID"})
		return
	}

	blog := bc.BlogUsecase.ViewBlogByID(id)
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// UpdateBlog modifies an existing blog post if the user is the author
func (bc *BlogController) UpdateBlog(c *gin.Context) {
	id := c.Param("id")

	var updatedBlog domain.Blog
	if err := c.ShouldBindJSON(&updatedBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve existing blog
	existingBlog := bc.BlogUsecase.ViewBlogByID(id)
	if existingBlog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	// Verify ownership
	userID := c.GetString("id")
	if existingBlog.AuthorID.Hex() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own blogs"})
		return
	}

	// Apply updates
	existingBlog.Title = updatedBlog.Title
	existingBlog.Content = updatedBlog.Content
	existingBlog.Tags = updatedBlog.Tags
	existingBlog.UpdatedAt = time.Now()

	if err := bc.BlogUsecase.UpdateBlog(id, existingBlog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingBlog)
}

// DeleteBlog removes a blog post if the user is the author or an admin
func (bc *BlogController) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	// Retrieve existing blog
	blog := bc.BlogUsecase.ViewBlogByID(id)
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	// Verify ownership or admin role
	userID := c.GetString("id")
	userRole := c.GetString("role")

	if blog.AuthorID.Hex() != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to delete this blog"})
		return
	}

	if err := bc.BlogUsecase.DeleteBlog(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

// ListBlogs retrieves paginated and filtered blog posts
func (bc *BlogController) ListBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Build filter from query parameters
	filter := domain.BlogFilter{
		Search:    c.Query("search"),
		Author:    c.Query("author"),
		Tag:       c.QueryArray("tag"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	// Parse optional date filters
	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate, _ = time.Parse(time.RFC3339, startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate, _ = time.Parse(time.RFC3339, endDate)
	}

	// Retrieve blogs
	blogs, total, err := bc.BlogUsecase.ListBlogs(page, limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blogs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": blogs,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}