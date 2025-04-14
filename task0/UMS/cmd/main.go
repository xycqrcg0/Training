package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ums/internal/config"
	"ums/internal/middlewares"
	"ums/internal/router"
)

func main() {
	config.InitConfig()

	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middlewares.CheckJWT())

	router.AuthRouter(r)
	router.CURDRouter(r)
	router.AdminRouter(r)

	r.Logger.Fatal(r.Start(config.Config.Port))
}
