package usecase

import (
	"fmt"
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogUseCase coordinates blog-related operations between repositories
type BlogUseCase struct {
	Repo            IBlogRepo
	InteractionRepo domain.InteractionRepository
	UserRepo        domain.UserRepository
}

// NewBlogUseCase initializes a new BlogUseCase with blog, interaction, and user repositories
func NewBlogUseCase(repo IBlogRepo, interactionRepo domain.InteractionRepository, urepo domain.UserRepository) *BlogUseCase {
	return &BlogUseCase{
		Repo:            repo,
		InteractionRepo: interactionRepo,
		UserRepo:        urepo,
	}
}

// StoreBlog creates a new blog post with default stats and timestamps
func (b *BlogUseCase) StoreBlog(blog *domain.Blog) error {
	// Set creation and update timestamps
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	// Initialize blog stats
	blog.Stats = domain.BlogStats{
		Views:    0,
		Likes:    0,
		Dislikes: 0,
		Comments: 0,
	}

	// Fetch author's username from UserRepo
	author := b.UserRepo.GetByID(blog.AuthorID)
	blog.AuthorName = author.Username

	// Store blog in repository
	err := b.Repo.StoreBlog(blog)
	if err != nil {
		fmt.Println("blog insertion failed")
		return err
	}

	fmt.Println("Inserted a blog")
	return nil
}

// ViewBlogByID retrieves a blog by its ID and asynchronously tracks the view
func (b *BlogUseCase) ViewBlogByID(blogID string) *domain.Blog {
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return &domain.Blog{}
	}

	// Track view in background
	go b.TrackView(blogID)

	// Retrieve blog from repository
	result := b.Repo.ViewBlogByID(id)
	return result
}

// GetBlogByAuthor retrieves blogs written by a specific author with pagination
func (b *BlogUseCase) GetBlogByAuthor(author string, page, limit int) ([]*domain.Blog, error) {
	// Validate pagination inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	skip := (page - 1) * 10

	return b.Repo.GetByAuthor(author, skip, limit)
}

// ViewBlogs retrieves all blogs from the repository
func (b *BlogUseCase) ViewBlogs() []domain.Blog {
	return b.Repo.RetriveAll()
}

// UpdateBlog modifies an existing blog post by its ID
func (b *BlogUseCase) UpdateBlog(blogID string, updatedBlog *domain.Blog) error {
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	// Update timestamp
	updatedBlog.UpdatedAt = time.Now()

	// Apply update via repository
	result := b.Repo.UpdateBlog(id, updatedBlog)
	return result
}

// DeleteBlog removes a blog post by its ID
func (b *BlogUseCase) DeleteBlog(blogID string) error {
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	return b.Repo.DeleteBlog(id)
}

// ListBlogs retrieves blogs with pagination and filtering options
func (b *BlogUseCase) ListBlogs(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
	// Set default sorting if not provided
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	return b.Repo.List(page, limit, filter)
}