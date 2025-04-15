package router

import "github.com/labstack/echo/v4"

func SetupRouter(r *echo.Echo) {
	AuthRouter(r)
	CURDRouter(r)
	AdminRouter(r)
}
