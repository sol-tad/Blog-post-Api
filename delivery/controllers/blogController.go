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

// CreateBlog godoc
// @Summary Create a new blog post
// @Description Create a new blog post by an authenticated user
// @Tags blogs
// @Accept json
// @Produce json
// @Param blog body domain.Blog true "Blog data"
// @Success 201 {object} domain.Blog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs [post]
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

	blog.AuthorID = objID
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	if err := bc.BlogUsecase.StoreBlog(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

// GetBlog godoc
// @Summary Get a blog post by ID
// @Description Retrieve a single blog post by its ID
// @Tags blogs
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} domain.Blog
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /blogs/{id} [get]
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

// UpdateBlog godoc
// @Summary Update a blog post
// @Description Update an existing blog post if you are the author
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param blog body domain.Blog true "Updated blog data"
// @Success 200 {object} domain.Blog
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs/{id} [put]
func (bc *BlogController) UpdateBlog(c *gin.Context) {
	id := c.Param("id")

	var updatedBlog domain.Blog
	if err := c.ShouldBindJSON(&updatedBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingBlog := bc.BlogUsecase.ViewBlogByID(id)
	if existingBlog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	userID := c.GetString("id")
	if existingBlog.AuthorID.Hex() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own blogs"})
		return
	}

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

// DeleteBlog godoc
// @Summary Delete a blog post
// @Description Delete a blog post if you are the author or an admin
// @Tags blogs
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs/{id} [delete]
func (bc *BlogController) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	blog := bc.BlogUsecase.ViewBlogByID(id)
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

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

// ListBlogs godoc
// @Summary List blog posts
// @Description Get paginated and filtered list of blog posts
// @Tags blogs
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Param author query string false "Author name"
// @Param tag query []string false "Tag filter"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /blogs [get]
func (bc *BlogController) ListBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := domain.BlogFilter{
		Search:    c.Query("search"),
		Author:    c.Query("author"),
		Tag:       c.QueryArray("tag"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate, _ = time.Parse(time.RFC3339, startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate, _ = time.Parse(time.RFC3339, endDate)
	}

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
