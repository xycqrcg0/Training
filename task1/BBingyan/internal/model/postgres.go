package model

import (
	"BBingyan/internal/config"
	"BBingyan/internal/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func newPostgres() {
	db, err := gorm.Open(postgres.Open(config.Config.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to connect to postgres")
	}
	DB = db

	//AutoMigrate
	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := DB.AutoMigrate(&Post{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := DB.AutoMigrate(&FollowShip{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := DB.AutoMigrate(&UserLikeShip{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}
	if err := DB.AutoMigrate(&PostLikeShip{}); err != nil {
		log.Fatalf("Fail to automigrate database")
	}

	log.Infof("Finish initializing postgres")
}
