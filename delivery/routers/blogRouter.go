package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)


func SetupBlogRoutes(router *gin.Engine) {

	// blogDbCollection:=config.BlogCollection
	// blogRepository:=repository.NewBlogRepository(blogDbCollection)
	// blogUsecase:=usecase.NewBlogUsecase(blogRepository)
	// blogController:=controllers.NewBlogController(blogUsecase)

	// blogRoutes:=router.Group("")

	{
		// blogRoutes.POST("",blogController.)
		// blogRoutes.POST("",blogController.)
	}
}