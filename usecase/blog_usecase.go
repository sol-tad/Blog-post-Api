package usecase

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/sol-tad/Blog-post-Api/domain"
)

type BlogUseCase struct {
	Repo IBlogRepo
}

func NewBlogUseCase(repo IBlogRepo) *BlogUseCase {
	return &BlogUseCase{
		Repo: repo,
	}
}

func (b *BlogUseCase) CreateBlog(blog *domain.Blog) {
	err := b.Repo.StoreBlog(blog)
	if err != nil {
		fmt.Println("blog insertion failed")
		return
	}
	fmt.Println("Inserted a blog")

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
	result := b.Repo.UpdateBlog(id, updatedBlog)
	return result
}

func (b *BlogUseCase) DeleteBlog (blogID string) error{
	
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	result := b.Repo.DeleteBlog(id)
	return result
}

