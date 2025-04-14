package models

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"ums/config"
	"ums/internal/global"
)

func InitPostgres() {
	db, err := gorm.Open(postgres.Open(config.Config.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to connect to database")
	}
	global.DB = db

	if err := global.DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Fail to init table")
	}

	//初始admin?存的时候就当普通用户吧
	_, err = GetUserByEmail("null")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if e := AddUser(&User{
				Name:     config.Config.Admin.Name,
				Email:    config.Config.Admin.Email,
				Password: config.Config.Admin.Password,
			}); e != nil {
				log.Fatalf("Fail to init admin")
			}
		} else {
			log.Fatalf("Fail to read users during initing admin")
		}
	}
}
