package usecase

import (
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
)

// CommentUsecase coordinates comment-related operations between repositories
type CommentUsecase struct {
	commentRepo domain.CommentRepository
	blogRepo    IBlogRepo
}

// NewCommentUsecase initializes a new CommentUsecase with comment and blog repositories
func NewCommentUsecase(
	commentRepo domain.CommentRepository,
	blogRepo IBlogRepo,
) *CommentUsecase {
	return &CommentUsecase{
		commentRepo: commentRepo,
		blogRepo:    blogRepo,
	}
}

// CreateComment adds a new comment and increments the blog's comment count
func (uc *CommentUsecase) CreateComment(comment *domain.Comment) error {
	// Set timestamps
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	// Increment comment count on the associated blog
	if err := uc.commentRepo.IncrementCommentCount(comment.BlogID.Hex()); err != nil {
		return err
	}

	// Store the comment
	return uc.commentRepo.Create(comment)
}

// GetCommentByID retrieves a comment by its ID
func (uc *CommentUsecase) GetCommentByID(id string) (*domain.Comment, error) {
	return uc.commentRepo.GetByID(id)
}

// GetCommentsByBlog retrieves paginated comments for a specific blog post
func (uc *CommentUsecase) GetCommentsByBlog(blogID string, page, limit int) ([]*domain.Comment, error) {
	return uc.commentRepo.GetByBlog(blogID, page, limit)
}

// UpdateComment modifies an existing comment
func (uc *CommentUsecase) UpdateComment(comment *domain.Comment) error {
	// Update timestamp
	comment.UpdatedAt = time.Now()

	// Apply update
	return uc.commentRepo.Update(comment)
}

// DeleteComment removes a comment and decrements the blog's comment count
func (uc *CommentUsecase) DeleteComment(id string) error {
	// Retrieve comment to access its BlogID
	comment, err := uc.commentRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Decrement comment count on the associated blog
	if err := uc.commentRepo.DecrementCommentCount(comment.BlogID.Hex()); err != nil {
		return err
	}

	// Delete the comment
	return uc.commentRepo.Delete(id)
}