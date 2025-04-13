package main

import (
	"BBingyan/config"
	"BBingyan/internal/router"
)

func main() {
	config.Config()
	r := router.SetupRouter()

	r.Logger.Fatalf(":9979")
}
