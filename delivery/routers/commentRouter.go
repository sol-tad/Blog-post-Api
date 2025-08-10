package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/sol-tad/Blog-post-Api/config"
    "github.com/sol-tad/Blog-post-Api/delivery/controllers"
    "github.com/sol-tad/Blog-post-Api/middlewares"
    "github.com/sol-tad/Blog-post-Api/repository"
    "github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupCommentRoutes registers routes related to blog comments.
func SetupCommentRoutes(router *gin.Engine) {
    // Initialize repositories
    commentRepo := repository.NewCommentRepository(config.CommentCollection)
    blogRepo := repository.NewBlogRepo(config.BlogCollection)

    // Initialize use case and controller
    commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo)
    commentController := controllers.NewCommentController(commentUsecase)

    // Group routes under /blogs/comments/:blog_id
    commentRoutes := router.Group("/blogs/comments/:blog_id")
    {
        // @Summary      Get all comments for a blog post
        // @Description  Retrieves all comments associated with the specified blog post ID.
        // @Tags         comments
        // @Param        blog_id   path      string  true  "Blog Post ID"
        // @Success      200  {object}  map[string]interface{}
        // @Failure      400  {object}  map[string]string  "Invalid blog ID"
        // @Router       /blogs/comments/{blog_id} [get]
        commentRoutes.GET("", commentController.GetComments)

        // Protected routes
        protected := commentRoutes.Group("")
        protected.Use(middlewares.AuthMiddleware())
        {
            // @Summary      Create a new comment
            // @Description  Adds a new comment to the specified blog post. Requires authentication.
            // @Tags         comments
            // @Param        blog_id   path      string  true  "Blog Post ID"
            // @Security     ApiKeyAuth
            // @Accept       json
            // @Produce      json
            // @Param        comment  body      domain.Comment  true  "Comment content"
            // @Success      201  {object}  domain.Comment
            // @Failure      400  {object}  map[string]string
            // @Failure      401  {object}  map[string]string
            // @Router       /blogs/comments/{blog_id} [post]
            protected.POST("", commentController.CreateComment)

            // @Summary      Update a comment
            // @Description  Updates an existing comment by comment_id for the specified blog post. Requires authentication.
            // @Tags         comments
            // @Param        blog_id     path      string  true  "Blog Post ID"
            // @Param        comment_id  path      string  true  "Comment ID"
            // @Security     ApiKeyAuth
            // @Accept       json
            // @Produce      json
            // @Param        comment  body      domain.Comment  true  "Updated comment content"
            // @Success      200  {object}  domain.Comment
            // @Failure      400  {object}  map[string]string
            // @Failure      401  {object}  map[string]string
            // @Failure      403  {object}  map[string]string
            // @Router       /blogs/comments/{blog_id}/{comment_id} [put]
            protected.PUT("/:comment_id", commentController.UpdateComment)

            // @Summary      Delete a comment
            // @Description  Deletes a comment by comment_id for the specified blog post. Requires authentication.
            // @Tags         comments
            // @Param        blog_id     path      string  true  "Blog Post ID"
            // @Param        comment_id  path      string  true  "Comment ID"
            // @Security     ApiKeyAuth
            // @Success      200  {object}  map[string]string  "Comment deleted successfully"
            // @Failure      400  {object}  map[string]string
            // @Failure      401  {object}  map[string]string
            // @Failure      403  {object}  map[string]string
            // @Router       /blogs/comments/{blog_id}/{comment_id} [delete]
            protected.DELETE("/:comment_id", commentController.DeleteComment)
        }
    }
}
