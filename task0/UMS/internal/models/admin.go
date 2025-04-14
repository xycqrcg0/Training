package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"ums/internal/config"
	"ums/internal/global"
)

type Admin struct {
	gorm.Model
	Name     string `gorm:"name"`
	Password string `gorm:"password"`
}

func InitAdmin() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(config.Config.Admin.Password), 12)
	if err != nil {
		return err
	}

	err = global.DB.Model(&Admin{}).Create(&Admin{
		Name:     config.Config.Admin.Name,
		Password: string(hashed),
	}).Error

	return err
}

func AddAdmin(newAdmin *Admin) error {
	err := global.DB.Model(&Admin{}).Create(newAdmin).Error
	return err
}

func GetAdminByName(name string) (*Admin, error) {
	var admin Admin
	err := global.DB.Model(&Admin{}).Where("name=?", name).First(&admin).Error
	return &admin, err
}
