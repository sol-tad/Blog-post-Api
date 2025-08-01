package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type BlogController struct{
	BlogUseCase *usecase.BlogUseCase
}

func NewBlogController (blog *usecase.BlogUseCase) *BlogController{
	return &BlogController{
		BlogUseCase: blog,
	}
}

type BlogDTO struct{
	ID 		primitive.ObjectID  `json:"_id"`
	Title   string 				`json:"title"`
	Content string 				`json:"content"`
	Tags    string 				`json:"tags"`
	Date    time.Time 			`json:"date"`
	User 	domain.User 		`json:"user"`
}
func (bc *BlogController) ChangeToBlog(blogDTO BlogDTO) *domain.Blog{
	var domBlog domain.Blog
	domBlog.Content = blogDTO.Content
	domBlog.Title = blogDTO.Title
	domBlog.Tags = blogDTO.Tags
	domBlog.Date = blogDTO.Date
	domBlog.User = blogDTO.User
	return &domBlog
}
func (bc *BlogController) CreateBlogController (ctx *gin.Context){
	var newBlog BlogDTO
	if err := ctx.ShouldBindJSON(&newBlog); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blog := bc.ChangeToBlog(newBlog)
	bc.BlogUseCase.CreateBlog(blog)
}

func (bc *BlogController) ViewBlogsController(ctx *gin.Context){
	blogs := bc.BlogUseCase.ViewBlogs()
	if len(blogs) == 0{
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "No blogs found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, blogs)
}

func (bc *BlogController) UpdateBlogController(ctx *gin.Context){
	id := ctx.Param("id")
	var updatedBlog BlogDTO
	if err := ctx.ShouldBindJSON(&updatedBlog); err != nil{
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainUpdatedBlog := bc.ChangeToBlog(updatedBlog)
	err := bc.BlogUseCase.UpdateBlog(id , domainUpdatedBlog)
		if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
}

func (bc *BlogController) DeleteBlogController(ctx *gin.Context) {
	id := ctx.Param("id")
	err := bc.BlogUseCase.DeleteBlog(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}