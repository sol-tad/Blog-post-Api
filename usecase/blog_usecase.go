package usecase

import (
	"fmt"
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUseCase struct {
	Repo IBlogRepo
	InteractionRepo domain.InteractionRepository
	UserRepo domain.UserRepository
}

func NewBlogUseCase(repo IBlogRepo, interactionRepo domain.InteractionRepository,urepo domain.UserRepository) *BlogUseCase {
	return &BlogUseCase{
		Repo: repo,
		InteractionRepo: interactionRepo,
		UserRepo: urepo,
	}
}


func (b *BlogUseCase) StoreBlog(blog *domain.Blog) error {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Stats = domain.BlogStats{
		Views:    0,
		Likes:    0,
		Dislikes: 0,
		Comments: 0,
	}
	author := b.UserRepo.GetByID(blog.AuthorID)
	blog.AuthorName = author.Username
	err := b.Repo.StoreBlog(blog)
	if err != nil {
		fmt.Println("blog insertion failed")
		return err
	}
	fmt.Println("Inserted a blog")
	return nil

}

func (b *BlogUseCase) ViewBlogByID(blogID string)*domain.Blog{
	id , err := primitive.ObjectIDFromHex(blogID)
	if err != nil{
		return &domain.Blog{}
	}
	
	go b.TrackView(blogID)
	result:= b.Repo.ViewBlogByID(id)

	return result
}
func (b *BlogUseCase) GetBlogByAuthor(author string, page, limit int) ([]*domain.Blog, error){
	 if page < 1 {page = 1}
	 if limit <1 || limit > 50 {limit = 10}
	 skip := (page - 1) * 10
	 return b.Repo.GetByAuthor(author, skip, limit)


}

func (b *BlogUseCase) ViewBlogs() []domain.Blog{
	return b.Repo.RetriveAll()
}

// there will be the updater and the blog's author
func (b *BlogUseCase) UpdateBlog(blogID string ,updatedBlog *domain.Blog) error{
	// b.Repo.CheckUserAuthority(blog , blog.User)
	id , err := primitive.ObjectIDFromHex(blogID)
	if err != nil{
		return err
	}
	updatedBlog.UpdatedAt = time.Now()
	result := b.Repo.UpdateBlog(id, updatedBlog)
	return result
}

func(b *BlogUseCase) DeleteBlog(blogID string) error {
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil{
		return err 
	}
	return b.Repo.DeleteBlog(id)

}
func (b *BlogUseCase) ListBlogs(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
	// Set default sort options
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}
	
	return b.Repo.List(page, limit, filter)
}