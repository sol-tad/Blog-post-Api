package routers

import (
	"github.com/gin-gonic/gin"
	// "github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	// "github.com/sol-tad/Blog-post-Api/repository"
	// "github.com/sol-tad/Blog-post-Api/usecase"
)


func SetupBlogRoutes(router *gin.Engine,  bc *controllers.BlogController ) {

	// endpoints: create, reterive, update,delete
	router.POST("/createblog" , bc.CreateBlogController)
	router.GET("/viewblogs", bc.ViewBlogsController)
	router.PUT("/updateblog" , bc.UpdateBlogController)
	router.DELETE("/deleteblog", bc.DeleteBlogController)
}