package model

import (
	"BBingyan/internal/config"
	"BBingyan/internal/global"
	"BBingyan/internal/log"
	"github.com/go-redis/redis"
)

func newRedis() {
	global.RedisDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		DB:       0,
		Password: "",
	})
	if _, err := global.RedisDB.Ping().Result(); err != nil {
		log.Fatalf("Fail to init redis")
	}
	log.Infof("Finish initializing redis")
}
