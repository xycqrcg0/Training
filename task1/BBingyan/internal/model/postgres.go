package model

import (
	"BBingyan/internal/config"
	"BBingyan/internal/global"
	"BBingyan/internal/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPostgres() {
	db, err := gorm.Open(postgres.Open(config.Config.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to connect to postgres")
	}
	global.DB = db

	//AutoMigrate

	log.Infof("Finish initializing postgres")
}
