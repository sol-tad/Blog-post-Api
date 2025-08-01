package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
)

func SetupRouter() *gin.Engine{
	router:=gin.Default()
	 // user routes
    SetupUserRoutes(router)

	
	// blog routes
	var bc *controllers.BlogController
	SetupBlogRoutes(router, bc)
	return router
}