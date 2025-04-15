package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ums/internal/config"
	"ums/internal/middlewares"
	"ums/internal/models"
	"ums/internal/router"
)

func main() {
	config.InitConfig()
	models.InitPostgres()

	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middlewares.CheckJWT())

	router.SetupRouter(r)

	r.Logger.Fatal(r.Start(config.Config.Port))
}
