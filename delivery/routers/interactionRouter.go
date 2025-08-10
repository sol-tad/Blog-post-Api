package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/sol-tad/Blog-post-Api/config"
    "github.com/sol-tad/Blog-post-Api/delivery/controllers"
    "github.com/sol-tad/Blog-post-Api/middlewares"
    "github.com/sol-tad/Blog-post-Api/repository"
    "github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupInteractionRoutes registers routes related to blog post interactions,
// such as liking and disliking posts. It initializes the interaction repository,
// use case, and controller, then sets up authenticated routes under the
// "/blogs/:id" path to handle user interactions on specific blog posts.
func SetupInteractionRoutes(router *gin.Engine) {
    // Initialize the interaction repository with blog and interaction collections
    interactionRepo := repository.NewInteractionRepository(
        config.BlogCollection,
        config.InteractionCollection,
    )

    // Create the interaction use case with the repository
    interactionUsecase := usecase.NewInteractionUsecase(interactionRepo)

    // Initialize the interaction controller with the use case
    interactionController := controllers.NewInteractionController(interactionUsecase)

    // Create a route group for routes under /blogs/:id, representing individual blog posts
    interactionRoutes := router.Group("/blogs/:id")

    // Protect all interaction routes with authentication middleware
    interactionRoutes.Use(middlewares.AuthMiddleware())

    {
        // @Summary      Like a blog post
        // @Description  Adds a like to the specified blog post for the authenticated user.
        // @Tags         interactions
        // @Param        id   path      string  true  "Blog Post ID"
        // @Security     BearerAuth
        // @Success      200  {object}  map[string]string  "Like added successfully"
        // @Failure      400  {object}  map[string]string  "Invalid request"
        // @Failure      401  {object}  map[string]string  "Unauthorized"
        // @Router       /blogs/{id}/like [put]
        interactionRoutes.PUT("/like", interactionController.LikeBlog)

        // @Summary      Unlike a blog post
        // @Description  Removes a like from the specified blog post for the authenticated user.
        // @Tags         interactions
        // @Param        id   path      string  true  "Blog Post ID"
        // @Security     BearerAuth
        // @Success      200  {object}  map[string]string  "Like removed successfully"
        // @Failure      400  {object}  map[string]string  "Invalid request"
        // @Failure      401  {object}  map[string]string  "Unauthorized"
        // @Router       /blogs/{id}/unlike [put]
        interactionRoutes.PUT("/unlike", interactionController.UnlikeBlog)

        // @Summary      Dislike a blog post
        // @Description  Adds a dislike to the specified blog post for the authenticated user.
        // @Tags         interactions
        // @Param        id   path      string  true  "Blog Post ID"
        // @Security     BearerAuth
        // @Success      200  {object}  map[string]string  "Dislike added successfully"
        // @Failure      400  {object}  map[string]string  "Invalid request"
        // @Failure      401  {object}  map[string]string  "Unauthorized"
        // @Router       /blogs/{id}/dislike [put]
        interactionRoutes.PUT("/dislike", interactionController.DislikeBlog)

        // @Summary      Remove a dislike from a blog post
        // @Description  Removes a previously added dislike from the specified blog post for the authenticated user.
        // @Tags         interactions
        // @Param        id   path      string  true  "Blog Post ID"
        // @Security     BearerAuth
        // @Success      200  {object}  map[string]string  "Dislike removed successfully"
        // @Failure      400  {object}  map[string]string  "Invalid request"
        // @Failure      401  {object}  map[string]string  "Unauthorized"
        // @Router       /blogs/{id}/undislike [put]
        interactionRoutes.PUT("/undislike", interactionController.UndoDislike)
    }
}
