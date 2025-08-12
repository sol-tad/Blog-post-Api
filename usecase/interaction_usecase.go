package usecase

import (
	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InteractionUsecase provides a dedicated layer for blog interaction logic
type InteractionUsecase struct {
	InteractionRepo domain.InteractionRepository
}

// NewInteractionUsecase initializes a new InteractionUsecase
func NewInteractionUsecase(InteractionRepo domain.InteractionRepository) *InteractionUsecase {
	return &InteractionUsecase{
		InteractionRepo: InteractionRepo,
	}
}

// TrackView asynchronously increments the view count for a blog post
func (b *BlogUseCase) TrackView(blogID string) {
	go b.InteractionRepo.IncrementViewCount(blogID)
}

// LikeBlog registers a like from a user for a blog post
func (b *BlogUseCase) LikeBlog(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.AddLike(blogID, userID)
}

// DislikeBlog registers a dislike from a user for a blog post
func (b *BlogUseCase) DislikeBlog(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.AddDislike(blogID, userID)
}

// RemoveLike unregisters a user's like from a blog post
func (b *BlogUseCase) RemoveLike(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.RemoveLike(blogID, userID)
}

// RemoveDislike unregisters a user's dislike from a blog post
func (b *BlogUseCase) RemoveDislike(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.RemoveDislike(blogID, userID)
}