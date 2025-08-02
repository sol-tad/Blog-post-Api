package usecase

import (
	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBlogRepo interface {
	StoreBlog(blog *domain.Blog) error
	RetriveAll() []domain.Blog
	ViewBlogByID(blogID primitive.ObjectID) *domain.Blog
	UpdateBlog(id primitive.ObjectID, updatedTask *domain.Blog) error
	DeleteBlog(id primitive.ObjectID) error

	GetByAuthor(author string, limt, page string) ([]*domain.Blog, error)
	List(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error)
	IncrementViewCount(blogID string) error
	IncrementLikeCount(blogID string) error
	IncrementDislikeCount(blogID string) error
	DecrementLikeCount(blogID string) error
	DecrementDislikeCount(blogID string) error
	IncrementCommentCount(blogID string) error
	DecrementCommentCount(blogID string) error 
}





	

