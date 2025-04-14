package router

import (
	"github.com/labstack/echo/v4"
	controller2 "ums/internal/controller"
)

func CURDRouter(r *echo.Echo) {
	curd := r.Group("/curd")
	{
		curd.GET("/info", controller2.GetInfo)
		curd.POST("/update", controller2.UpdateUser)
		curd.POST("/delete", controller2.DeleteUser)
	}
}
