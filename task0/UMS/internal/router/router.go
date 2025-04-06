package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	controller2 "ums/internal/controller"
	"ums/internal/middlewares"
)

func SetupRouter() *echo.Echo {
	r := echo.New()

	r.Use(middleware.Logger())

	auth := r.Group("/auth")
	{
		auth.POST("/login", controller2.Login)
		auth.POST("/register", controller2.Register)
	}

	curd := r.Group("/curd")
	curd.Use(middlewares.CheckJWT())
	{
		curd.GET("/info", controller2.GetInfo)
		curd.POST("/update", controller2.UpdateUser)
		curd.POST("/delete", controller2.DeleteUser)
	}

	return r
}
