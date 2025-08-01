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
	ID 		primitive.ObjectID  `json:"id"`
	Title   string 				`json:"title"`
	Content string 				`json:"content"`
	Tags    []string 			`json:"tags"`
	AuthorID primitive.ObjectID `json:"author_id"`
	AuthorName string           `json:"author_name"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Stats     domain.BlogStats  `json:"stats"`
}
func (bc *BlogController) ChangeToBlog(blogDTO BlogDTO) *domain.Blog{
	var domBlog domain.Blog
	domBlog.Content = blogDTO.Content
	domBlog.ID = blogDTO.ID
	domBlog.Title = blogDTO.Title
	domBlog.Tags = blogDTO.Tags
	domBlog.CreatedAt = blogDTO.CreatedAt
	domBlog.UpdatedAt = blogDTO.UpdatedAt
	domBlog.AuthorID = blogDTO.AuthorID
	domBlog.AuthorName = blogDTO.AuthorName
	domBlog.Stats = blogDTO.Stats
	
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
	ctx.JSON(http.StatusOK, gin.H{"message": "blog created successfully"})
	
}

func (bc *BlogController) ViewBlogsController(ctx *gin.Context){
	blogs := bc.BlogUseCase.ViewBlogs()
	if len(blogs) == 0{
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "No blogs found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, blogs)
}
func (bc *BlogController) ViewBlogByIDController(ctx *gin.Context){
	id := ctx.Param("id")
	blog := bc.BlogUseCase.ViewBlogByID(id)
	ctx.IndentedJSON(200 , blog)
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
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}



// add the controllers for your additional endpoints here.

