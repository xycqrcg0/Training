package model

import (
	"BBingyan/internal/global"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"email;not null"`
	Name      string `gorm:"name;not null"`
	Password  string `gorm:"password;not null"`
	Signature string `gorm:"signature;default:你好，世界"`
}

func AddUser(newUser *User) error {
	err := global.DB.Model(&User{}).Create(newUser).Error
	return err
}

func DeleteUser(email string) error {
	err := global.DB.Model(&User{}).Where("email=?", email).Delete(&User{}).Error
	return err
}

func UpdateUser(user *User) error {
	err := global.DB.Model(&User{}).Where("email=?", user.Email).Updates(user).Error
	return err
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := global.DB.Model(&User{}).Where("email=?", email).First(user).Error
	return user, err
}

func GetAllUsersInfo() ([]User, error) {
	users := make([]User, 0)
	err := global.DB.Model(&User{}).Find(&users).Error

	return users, err
}
