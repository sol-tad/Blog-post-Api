package usecase

import (
	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type InteractionUsecase struct {
	InteractionRepo domain.InteractionRepository
}

func NewInteractionUsecase(InteractionRepo domain.InteractionRepository) *InteractionUsecase{
	return &InteractionUsecase{
		InteractionRepo: InteractionRepo,
	}
}



// TRACK THE BLOG VIEW 
func (b *BlogUseCase) TrackView(blogID string){
 
	go b.InteractionRepo.IncrementViewCount(blogID)
}

// LIKE AND DISLIKE FUNCTIONALITIES 
func (b *BlogUseCase) LikeBlog(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.AddLike(blogID, userID)
}

func (b *BlogUseCase) DislikeBlog(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.AddDislike(blogID, userID)

}
func (b *BlogUseCase) RemoveLike(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.RemoveLike(blogID, userID)
}

func (b *BlogUseCase) RemoveDislike(blogID string, userID primitive.ObjectID) error {
	return b.InteractionRepo.RemoveDislike(blogID, userID)
}














