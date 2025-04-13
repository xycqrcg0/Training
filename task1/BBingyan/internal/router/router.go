package router

import (
	"BBingyan/internal/controller"
	"BBingyan/internal/middleware"
	"github.com/labstack/echo/v4"
)

func SetupRouter() *echo.Echo {
	r := echo.New()

	auth := r.Group("/auth")
	{
		auth.POST("/register/code", controller.RegisterForCode)
		auth.POST("/register", controller.Register)
		auth.POST("/login/code", controller.LoginForCode)
		auth.POST("/login", controller.Login)
	}

	r.POST("/info", controller.UpdateInfo, middleware.CheckJWT())

	return r
}
