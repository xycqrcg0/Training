package main

import (
	"BBingyan/internal/config"
	"BBingyan/internal/router"
)

func main() {
	config.InitConfig()
	r := router.SetupRouter()

	r.Logger.Fatalf(config.Config.Port)
}
