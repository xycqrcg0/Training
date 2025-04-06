package main

import (
	"ums/config"
	"ums/internal/global"
	"ums/internal/router"
)

func main() {
	config.InitConfig()

	r := router.SetupRouter()

	r.Logger.Fatal(r.Start(global.Configs.Port))
}
