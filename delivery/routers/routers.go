package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	router:=gin.Default()
	
	 // user routes
    SetupUserRoutes(router)

	//ai router
	SetupAIRoutes(router)

	//like router
	SetupInteractionRoutes(router)
	
	//oauth routes
	SetupOAuthRouter(router)

	//comment router
	// SetupCommentRoutes(router)

	// blog routes
	SetupBlogRoutes(router)
	
	return router
}