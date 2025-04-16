package model

import (
	"BBingyan/internal/global"
	"gorm.io/gorm"
)

type FollowShip struct {
	gorm.Model
	User         string `gorm:"user"`
	FollowedUser string `gorm:"followed_user"`
}

func FollowUser(user string, followedUser string) error {
	err := global.DB.Model(&FollowShip{}).Create(&FollowShip{
		User:         user,
		FollowedUser: followedUser,
	}).Error
	return err
}

func UnfollowUser(user string, followedUser string) error {
	err := global.DB.Model(&FollowShip{}).Where("user=? AND followed_user=?", user, followedUser).Error
	return err
}

func GetFollows(user string) ([]FollowShip, error) {
	follows := make([]FollowShip, 0)
	err := global.DB.Model(&FollowShip{}).Where("user=?", user).Find(&follows).Error
	return follows, err
}
