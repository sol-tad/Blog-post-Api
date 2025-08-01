package usecase

import (
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
)

type BlogUsecase struct {
	blogRepo domain.BlogRepository
}

func NewBlogUsecase(blogRepo domain.BlogRepository) *BlogUsecase {
	return &BlogUsecase{
		blogRepo: blogRepo,
	}
}

func (uc *BlogUsecase) CreateBlog(blog *domain.Blog) error {
	// Set creation and update times
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	
	// Initialize stats
	blog.Stats = domain.BlogStats{
		Views:    0,
		Likes:    0,
		Dislikes: 0,
		Comments: 0,
	}
	
	return uc.blogRepo.Create(blog)
}

func (uc *BlogUsecase) GetBlogByID(id string) (*domain.Blog, error) {
	// Track view asynchronously
	go uc.blogRepo.IncrementViewCount(id)
	
	return uc.blogRepo.GetByID(id)
}

func (uc *BlogUsecase) UpdateBlog(blog *domain.Blog) error {
	blog.UpdatedAt = time.Now()
	return uc.blogRepo.Update(blog)
}

func (uc *BlogUsecase) DeleteBlog(id string) error {
	return uc.blogRepo.Delete(id)
}

func (uc *BlogUsecase) ListBlogs(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
	// Set default sort options
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}
	
	return uc.blogRepo.List(page, limit, filter)
}

func (uc *BlogUsecase) LikeBlog(blogID string) error {
	return uc.blogRepo.IncrementLikeCount(blogID)
}

func (uc *BlogUsecase) UnlikeBlog(blogID string) error {
	return uc.blogRepo.DecrementLikeCount(blogID)
}

func (uc *BlogUsecase) DislikeBlog(blogID string) error {
	return uc.blogRepo.IncrementDislikeCount(blogID)
}

func (uc *BlogUsecase) UndoDislikeBlog(blogID string) error {
	return uc.blogRepo.DecrementDislikeCount(blogID)
}