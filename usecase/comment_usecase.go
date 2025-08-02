package usecase

import (
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
)

type CommentUsecase struct {
    commentRepo domain.CommentRepository
    blogRepo    IBlogRepo
}

func NewCommentUsecase(
    commentRepo domain.CommentRepository,
    blogRepo IBlogRepo,
) *CommentUsecase {
    return &CommentUsecase{
        commentRepo: commentRepo,
        blogRepo:    blogRepo,
    }
}

func (uc *CommentUsecase) CreateComment(comment *domain.Comment) error {
    comment.CreatedAt = time.Now()
    comment.UpdatedAt = time.Now()
    
    // Increment blog comment count
    if err := uc.blogRepo.IncrementCommentCount(comment.BlogID.Hex()); err != nil {
        return err
    }
    
    return uc.commentRepo.Create(comment)
}

func (uc *CommentUsecase) GetCommentByID(id string) (*domain.Comment, error) {
    return uc.commentRepo.GetByID(id)
}

func (uc *CommentUsecase) GetCommentsByBlog(blogID string, page, limit int) ([]*domain.Comment, error) {
    return uc.commentRepo.GetByBlog(blogID, page, limit)
}

func (uc *CommentUsecase) UpdateComment(comment *domain.Comment) error {
    comment.UpdatedAt = time.Now()
    return uc.commentRepo.Update(comment)
}

func (uc *CommentUsecase) DeleteComment(id string) error {
    comment, err := uc.commentRepo.GetByID(id)
    if err != nil {
        return err
    }
    // Decrement blog comment count
    if err := uc.blogRepo.DecrementCommentCount(comment.BlogID.Hex()); err != nil {
        return err
    }
    
    return uc.commentRepo.Delete(id)
}