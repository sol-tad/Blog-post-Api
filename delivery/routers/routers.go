package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	router:=gin.Default()
	
	 // user routes
    SetupUserRoutes(router)

	SetupInteractionRoutes(router)
	
	//oauth routes
	SetupOAuthRouter(router)


	// blog routes
	SetupBlogRoutes(router)

	//ai routes
	SetupAiRoutes(router)
	return router
}