package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)


func SetupOAuthRouter(router *gin.Engine) {

	userDbCollection:=config.UserCollection

	userRepository:=repository.NewUserRepository(userDbCollection)

	oauthUsecase := usecase.NewOAuthUsecase(userRepository)
	oauthController:=controllers.NewOAuthController(oauthUsecase)

	oauthRoutes:=router.Group("")

	{
		
		oauthRoutes.GET("/oauth/google/login",oauthController.Login)
		oauthRoutes.GET("/oauth/google/callback",oauthController.Callback)

	}
}



