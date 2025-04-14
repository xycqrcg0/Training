package models

import (
	"gorm.io/gorm"
	"time"
	"ums/internal/global"
)

type User struct {
	Id        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func AddUser(newUser *User) error {
	err := global.DB.Model(&User{}).Create(newUser).Error
	return err
}

func DeleteUser(email string) error {
	err := global.DB.Model(&User{}).Where("email=?", email).Error
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
