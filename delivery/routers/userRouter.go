package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupUserRoutes initializes all user-related routes on the provided Gin router.
// It sets up the user repository, use case, and controller, then registers routes
// for user registration, login, profile management, password reset, OTP verification,
// token refresh, logout, and user role promotion/demotion.
//
// Some routes are protected by the AuthMiddleware to ensure only authenticated users can access them.
func SetupUserRoutes(router *gin.Engine) {
	userDbCollection := config.UserCollection

	userRepository := repository.NewUserRepository(userDbCollection)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controllers.NewUserController(userUsecase)

	// Public routes: no authentication required
	public := router.Group("")
	{
		public.POST("/register", userController.Register)
		public.POST("/verify-otp", userController.VerifyOTP)
		public.POST("/login", userController.Login)
		public.POST("/refresh", userController.RefreshTokenController)
		public.POST("/forgot-password", userController.SendResetOTP)
		public.POST("/reset-password", userController.ResetPassword)
	}

	// Protected routes: require user authentication
	protected := router.Group("")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/logout", userController.Logout)
		protected.PUT("/profile", userController.UpdateProfile)

		// Admin routes - for promoting/demoting users
		protected.POST("/user/:id/promote", userController.PromoteUser)
		protected.POST("/user/:id/demote", userController.DemoteUser)
	}
}
