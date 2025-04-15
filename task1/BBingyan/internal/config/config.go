package config

import (
	"BBingyan/internal/global"
	"BBingyan/internal/log"
	"encoding/json"
	"github.com/joho/godotenv"
	"os"
)

type PostgresConfig struct {
	Dsn string `json:"dsn"`
}

type RedisConfig struct {
	Addr string `json:"addr"`
}

type AdminConfig struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type JwtConfig struct {
	Key     string   `json:"key"`
	Exp     int      `json:"exp"`
	Skipper []string `json:"skipper"`
}

type StructConfig struct {
	Port     string         `json:"port"`
	Postgres PostgresConfig `json:"postgres"`
	Redis    RedisConfig    `json:"redis"`
	Admin    AdminConfig    `json:"admin"`
	JWT      JwtConfig      `json:"jwt"`
}

var Config StructConfig

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		global.Errors.Fatalf("fail to load .env file")
	}
	global.AuthorizationCode = os.Getenv("AUTH_CODE")
	global.Key = []byte(os.Getenv("KEY"))

	file, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalf("Fail to read from cinfig.json")
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalf("Fail to unmarshal config.json")
	}
	log.Infof("finish initializing config")
}
