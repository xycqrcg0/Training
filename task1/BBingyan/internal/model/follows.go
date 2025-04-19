package model

import (
	"BBingyan/internal/global"
	"gorm.io/gorm"
)

type FollowShip struct {
	gorm.Model
	User         string `gorm:"user"`
	FollowedUser string `gorm:"followed_user"`
	Info         User   `gorm:"foreignKey:User;references:Email"`
	FollowedInfo User   `gorm:"foreignKey:FollowedUser;references:Email"`
}

//这里的改动会影响到user的follows字段，在这里开事务并修改user的follows字段应该没问题吧

func FollowUser(user string, followedUser string) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&FollowShip{}).Create(&FollowShip{
			User:         user,
			FollowedUser: followedUser,
		})
		tx.Model(&User{}).Where("email=?", followedUser).Update("follows", gorm.Expr("follows+1"))
		return nil
	})
	return err
}

func UnfollowUser(user string, followedUser string) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&FollowShip{}).Where("user=? AND followed_user=?", user, followedUser).Delete(&FollowShip{})
		tx.Model(&User{}).Where("email=?", followedUser).Update("follows", gorm.Expr("follows-1"))
		return nil
	})
	return err
}

func GetAllFollows(user string, page int, pageSize int) ([]FollowShip, error) {
	follows := make([]FollowShip, 0)

	err := global.DB.Model(&FollowShip{}).Preload("FollowedInfo").Where("user=?", user).
		Limit(pageSize).Offset(page).Find(&follows).Error

	return follows, err
}

func GetAllFans(user string, page int, pageSize int) ([]FollowShip, error) {
	fans := make([]FollowShip, 0)

	err := global.DB.Model(&FollowShip{}).Preload("Info").Where("followed_user=?", user).
		Limit(pageSize).Offset(page).Find(&fans).Error

	return fans, err
}
