package router

import (
	"github.com/labstack/echo/v4"
	controller2 "ums/internal/controller"
)

func AuthRouter(r *echo.Echo) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", controller2.Login)
		auth.POST("/register", controller2.Register)
	}
}
