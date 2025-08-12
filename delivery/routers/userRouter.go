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
	// Get the MongoDB user collection from config
	userDbCollection := config.UserCollection

	// Initialize the user repository with the user collection
	userRepository := repository.NewUserRepository(userDbCollection)

	// Create the user use case instance with the repository
	userUsecase := usecase.NewUserUsecase(userRepository)

	// Initialize the user controller with the use case
	userController := controllers.NewUserController(userUsecase)

	// Create a route group without a specific prefix (root level)
	userRoutes := router.Group("")

	{
		// Public routes
		userRoutes.POST("/register", userController.Register)           // Register a new user
		userRoutes.POST("/verify-otp", userController.VerifyOTP)        // Verify OTP sent to user
		userRoutes.POST("/login", userController.Login)                 // User login
		userRoutes.POST("/refresh", userController.RefreshTokenController) // Refresh JWT tokens
		userRoutes.POST("/forgot-password", userController.SendResetOTP)  // Request password reset OTP
		userRoutes.POST("/reset-password", userController.ResetPassword)  // Reset password with OTP

		// Protected routes (require authentication)
		userRoutes.POST("/logout", middlewares.AuthMiddleware(), userController.Logout) // User logout
		userRoutes.PUT("/profile", middlewares.AuthMiddleware(), userController.UpdateProfile) // Update user profile

		// Admin or privileged user routes (require authentication)
		userRoutes.POST("/user/:id/promote", middlewares.AuthMiddleware(), userController.PromoteUser) // Promote user role
		userRoutes.POST("/user/:id/demote", middlewares.AuthMiddleware(), userController.DemoteUser)   // Demote user role
	}
}
