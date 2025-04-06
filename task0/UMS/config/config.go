package config

import (
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"ums/internal/global"
)

//var dsn = "host=localhost user=postgres password=123456 dbname=ums port=3456 sslmode=disable"

func InitConfig() {
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatalf("fail to open config.json")
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&global.Configs); err != nil {
		log.Fatalf("fail to read config.json")
	}

	file.Close()

	db, err := gorm.Open(postgres.Open(global.Configs.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败")
	}
	global.DB = db

}
