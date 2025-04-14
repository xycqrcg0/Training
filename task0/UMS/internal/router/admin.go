package router

import (
	"github.com/labstack/echo/v4"
	"ums/internal/controller"
)

func AdminRouter(r *echo.Echo) {
	admin := r.Group("/admin")
	{
		admin.POST("/login", controller.AdminLogin)
		admin.POST("/new-admin", controller.AddNewAdmin)
		admin.GET("/curd/users/:email", controller.GetUserInfo)
		admin.GET("/curd/users", controller.GetAllUsers)
		admin.GET("/curd/drop/:email", controller.DropUserByEmail)
	}
}
