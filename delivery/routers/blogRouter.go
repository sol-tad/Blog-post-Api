package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/middlewares"
	"github.com/sol-tad/Blog-post-Api/repository"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

// SetupBlogRoutes registers routes related to blog posts.
func SetupBlogRoutes(router *gin.Engine) {
	blogRepo := repository.NewBlogRepo(config.BlogCollection)
	userRepo := repository.NewUserRepository(config.UserCollection)
	interactionRepo := repository.NewInteractionRepository(
		config.BlogCollection,
		config.InteractionCollection,
	)

	blogUsecase := usecase.NewBlogUseCase(blogRepo, interactionRepo, userRepo)
	blogController := controllers.NewBlogController(blogUsecase)

	blogRoutes := router.Group("/blogs")
	{
		// @Summary      List all blogs
		// @Description  Retrieves a list of all blog posts.
		// @Tags         blogs
		// @Produce      json
		// @Success      200  {array}   map[string]interface{}
		// @Failure      500  {object}  map[string]string
		// @Router       /blogs [get]
		blogRoutes.GET("", blogController.ListBlogs)

		// @Summary      Get blog by ID
		// @Description  Retrieves a single blog post by its ID.
		// @Tags         blogs
		// @Produce      json
		// @Param        id   path      string  true  "Blog ID"
		// @Success      200  {object}  map[string]interface{}
		// @Failure      404  {object}  map[string]string
		// @Router       /blogs/{id} [get]
		blogRoutes.GET("/:id", blogController.GetBlog)

		protected := blogRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			// @Summary      Create a new blog
			// @Description  Creates a new blog post. Requires authentication.
			// @Tags         blogs
			// @Security     BearerAuth
			// @Accept       json
			// @Produce      json
			// @Param        blog  body      map[string]interface{}  true  "Blog data"
			// @Success      201  {object}  map[string]interface{}
			// @Failure      400  {object}  map[string]string
			// @Failure      401  {object}  map[string]string
			// @Router       /blogs/create [post]
			protected.POST("/create", blogController.CreateBlog)

			// @Summary      Update a blog
			// @Description  Updates an existing blog post by ID. Requires authentication.
			// @Tags         blogs
			// @Security     BearerAuth
			// @Accept       json
			// @Produce      json
			// @Param        id    path      string  true  "Blog ID"
			// @Param        blog  body      map[string]interface{}  true  "Updated blog data"
			// @Success      200  {object}  map[string]interface{}
			// @Failure      400  {object}  map[string]string
			// @Failure      401  {object}  map[string]string
			// @Failure      404  {object}  map[string]string
			// @Router       /blogs/{id} [put]
			protected.PUT("/:id", blogController.UpdateBlog)

			// @Summary      Delete a blog
			// @Description  Deletes a blog post by ID. Requires authentication.
			// @Tags         blogs
			// @Security     BearerAuth
			// @Param        id   path      string  true  "Blog ID"
			// @Success      200  {object}  map[string]string  "Blog deleted successfully"
			// @Failure      401  {object}  map[string]string
			// @Failure      404  {object}  map[string]string
			// @Router       /blogs/{id} [delete]
			protected.DELETE("/:id", blogController.DeleteBlog)
		}
	}
}
