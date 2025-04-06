package main

import (
	"BBingyan/config"
	"BBingyan/internal/util"
	"log"
)

func main() {
	config.Config()
	err := util.SendAuthCode("3299511912@qq.com", "1024")
	if err != nil {
		log.Fatalf("%v", err)
	}
}
