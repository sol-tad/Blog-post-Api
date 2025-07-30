package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)


func SetupUserRoutes(router *gin.Engine) {

	userDbCollection:=config.UserCollection

	userRepository:=repository.NewUserRepository(userDbCollection)
	userUsecase:=usecase.NewUserUsecase(userRepository)
	userController:=controllers.NewUserController(userUsecase)

	userRoutes:=router.Group("")

	{
		userRoutes.POST("/login",userController.Login)
		userRoutes.POST("/register",userController.Register)
	}
}



