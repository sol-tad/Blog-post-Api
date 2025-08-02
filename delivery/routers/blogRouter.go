package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

func SetupBlogRoutes(router *gin.Engine) {
	blogDbCollection := config.BlogCollection

	blogRepository := repository.NewBlogRepo(blogDbCollection)
	blogUseCase := usecase.NewBlogUseCase(blogRepository)
	blogController := controllers.NewBlogController(blogUseCase)

	// endpoints: create, reterive, update,delete
	router.POST("/createblog", blogController.CreateBlogController)
	router.GET("/viewblogs", blogController.ViewBlogsController)
	router.PUT("/updateblog/:id", blogController.UpdateBlogController)
	router.DELETE("/deleteblog/:id", blogController.DeleteBlogController)
	router.GET("/viewblogbyid/:id", blogController.ViewBlogByIDController)


// for your work, if there is a router needed, you can add it here.
}
