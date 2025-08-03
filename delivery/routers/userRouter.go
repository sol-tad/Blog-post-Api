package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
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
		
		userRoutes.POST("/register",userController.Register)
		userRoutes.POST("/verify-otp",userController.VerifyOTP)
		userRoutes.POST("/login",userController.Login)
		userRoutes.POST("/refresh",userController.RefreshTokenController)
		userRoutes.POST("/logout",middlewares.AuthMiddleware(),userController.Logout)
		
		userRoutes.POST("/forgot-password", userController.SendResetOTP)
		userRoutes.POST("/reset-password", userController.ResetPassword)

		userRoutes.PUT("/profile", middlewares.AuthMiddleware(), userController.UpdateProfile)


		userRoutes.POST("/user/:id/promote", middlewares.AuthMiddleware(), userController.PromoteUser)
		userRoutes.POST("/user/:id/demote", middlewares.AuthMiddleware(), userController.DemoteUser)

	}
}



