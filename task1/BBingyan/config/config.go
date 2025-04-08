package config

import (
	"BBingyan/internal/global"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
)

func Config() {
	//配errors.log
	w1 := os.Stdout
	w2, err := os.OpenFile("./internal/log/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("fail to open errors.log file")
	}
	global.Errors = logrus.New()
	global.Errors.SetOutput(io.MultiWriter(w1, w2))
	global.Errors.SetFormatter(&logrus.JSONFormatter{})
	global.Errors.SetReportCaller(true)

	//配info.log
	w3, err := os.OpenFile("./internal/log/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("fail to open info.log file")
	}
	global.Infos = logrus.New()
	global.Infos.SetOutput(io.MultiWriter(w1, w3))
	global.Errors.SetFormatter(&logrus.JSONFormatter{})

	if err := godotenv.Load(); err != nil {
		global.Errors.Fatalf("fail to load .env file")
	}
	global.AuthorizationCode = os.Getenv("AUTH_CODE")

	global.Key = []byte(os.Getenv("KEY"))

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		global.Errors.Fatalf("Fail to open postgres")
	}
	global.DB = db

	redisAddr := os.Getenv("REDIS_ADDR")
	global.RedisDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       0,
		Password: "",
	})
	if _, err := global.RedisDB.Ping().Result(); err != nil {
		global.Errors.Fatalf("Fail to open redis")
	}

	global.Infos.Info("finish initialization")
}
