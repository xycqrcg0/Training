package router

import (
	"BBingyan/internal/controller"
	"github.com/labstack/echo/v4"
)

func PostRouter(r *echo.Echo) {
	post := r.Group("/posts")
	{
		post.POST("", controller.AddPost)
		post.DELETE("/:id", controller.DeletePost)
		post.GET("/:email", controller.GetPostByEmail)
		post.GET("", controller.GetPostByTag)
	}
}
