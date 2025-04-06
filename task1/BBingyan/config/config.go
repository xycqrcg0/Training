package config

import (
	"BBingyan/internal/global"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func Config() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("fail to load .env file")
	}
	global.AuthorizationCode = os.Getenv("AUTH_CODE")

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to open postgres")
	}
	global.DB = db

	redisAddr := os.Getenv("REDIS_ADDR")
	global.RedisDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       0,
		Password: "",
	})
	if _, err := global.RedisDB.Ping().Result(); err != nil {
		log.Fatalf("Fail to open redis")
	}

}
