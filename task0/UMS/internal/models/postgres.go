package models

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"ums/internal/config"
	"ums/internal/global"
)

func InitPostgres() {
	db, err := gorm.Open(postgres.Open(config.Config.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to connect to database")
	}
	global.DB = db

	if err := global.DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Fail to init users table")
	}

	if err := global.DB.AutoMigrate(&Admin{}); err != nil {
		log.Fatalf("Fail to init admins table")
	}

	//初始admin?存的时候就当普通用户吧
	_, err = GetAdminByName(config.Config.Admin.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if e := InitAdmin(); e != nil {
				log.Fatalf("Fail to init admin")
			}
		} else {
			log.Fatalf("Fail to read users during initing admin")
		}
	}
}
