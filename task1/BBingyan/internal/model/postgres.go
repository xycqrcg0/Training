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
	if err := global.DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := global.DB.AutoMigrate(&Post{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := global.DB.AutoMigrate(&FollowShip{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := global.DB.AutoMigrate(&UserLikeShip{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := global.DB.AutoMigrate(&PostLikeShip{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}

	log.Infof("Finish initializing postgres")
}
