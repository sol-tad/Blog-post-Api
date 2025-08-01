package usecase

import (
	"fmt"
	"time"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
		blog.Stats = domain.BlogStats{
		Views:    0,
		Likes:    0,
		Dislikes: 0,
		Comments: 0,
	}
	err := b.Repo.StoreBlog(blog)
	if err != nil {
		fmt.Println("blog insertion failed")
		return
	}
	fmt.Println("Inserted a blog")

}

func (b *BlogUseCase) ViewBlogByID(blogID string)*domain.Blog{
	id , err := primitive.ObjectIDFromHex(blogID)
	if err != nil{
		return &domain.Blog{}
	}
	result:= b.Repo.ViewBlogByID(id)

	return result
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

func (b *BlogUseCase) DeleteBlog (blogID string) error{
	
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	result := b.Repo.DeleteBlog(id)
	return result
}




// I dont understand the things below, I mean how you expected them to work.
//so I just copied it here and modified some things. Your work would be implementing 
//the repo and the interface.




// func (b *BlogUseCase) ListBlogs(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
// 	// Set default sort options
// 	if filter.SortBy == "" {
// 		filter.SortBy = "created_at"
// 	}
// 	if filter.SortOrder == "" {
// 		filter.SortOrder = "desc"
// 	}
	
// 	return b.Repo.List(page, limit, filter)
// }

// func (b *BlogUseCase) LikeBlog(blogID string) error {
// 	return b.Repo.IncrementLikeCount(blogID)
// }

// func (b *BlogUseCase) UnlikeBlog(blogID string) error {
// 	return b.Repo.DecrementLikeCount(blogID)
// }

// func (b *BlogUseCase) DislikeBlog(blogID string) error {
// 	return b.Repo.IncrementDislikeCount(blogID)
// }

// func (b *BlogUseCase) UndoDislikeBlog(blogID string) error {
// 	return b.Repo.DecrementDislikeCount(blogID)
// }