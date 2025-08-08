package usecase

import (
	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IBlogRepo defines the contract for blog-related data operations
type IBlogRepo interface {
	// StoreBlog saves a new blog post to the database
	StoreBlog(blog *domain.Blog) error

	// RetriveAll returns all blog posts (useful for admin or public feed)
	RetriveAll() []domain.Blog

	// ViewBlogByID retrieves a single blog post by its ObjectID
	ViewBlogByID(blogID primitive.ObjectID) *domain.Blog

	// UpdateBlog modifies an existing blog post
	UpdateBlog(id primitive.ObjectID, updatedTask *domain.Blog) error

	// DeleteBlog removes a blog post by its ObjectID
	DeleteBlog(id primitive.ObjectID) error

	// GetByAuthor retrieves blog posts written by a specific author with pagination
	GetByAuthor(author string, skip, limit int) ([]*domain.Blog, error)

	// List returns filtered and paginated blog posts along with total count
	List(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error)
}