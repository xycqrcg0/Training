package config

import (
	"encoding/json"
	"log"
	"os"
	"ums/internal/models"
)

type PostgresConfig struct {
	Dsn string `json:"dsn"`
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
	Admin    AdminConfig    `json:"admin"`
	JWT      JwtConfig      `json:"jwt"`
}

var Config StructConfig

func InitConfig() {
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatalf("fail to open config.json")
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&Config); err != nil {
		log.Fatalf("fail to read config.json")
	}

	file.Close()

	models.InitPostgres()
}
