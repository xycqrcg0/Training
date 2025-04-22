package model

import (
	"BBingyan/internal/global"
	"errors"
	"gorm.io/gorm"
)

//id是能暴露的吗

type Post struct {
	gorm.Model
	Author  string `gorm:"author;not null"`
	Title   string `gorm:"title;not null"`
	Tag     string `gorm:"tag"`
	Content string `gorm:"content;not null"`
	Likes   int    `gorm:"likes"`
	Replies int    `gorm:"replies"`
	User    User   `gorm:"user;foreignKey:Author;references:Email"`
}

func AddPost(newPost *Post) error {
	err := global.DB.Model(&Post{}).Create(newPost).Error
	return err
}

func DeletePostById(user string, id int) error {
	result := global.DB.Model(&Post{}).Where("id=? AND author=?", id, user).Delete(&Post{})
	if result.RowsAffected == 0 {
		return errors.New("none")
	}
	return result.Error
}

func GetPostById(id int) (*Post, error) {
	post := &Post{}
	err := global.DB.Model(&Post{}).Preload("User").Where("id=?", id).First(post).Error
	return post, err
}

func GetPostLikes(id int) (int, error) {
	var likes int
	err := global.DB.Model(&Post{}).Select("likes").Where("id=?", id).First(&likes).Error
	return likes, err
}

func GetPostsByEmail(email string, page int, pageSize int) ([]Post, error) {
	posts := make([]Post, 0)
	err := global.DB.Model(&Post{}).Preload("User").Where("author=?", email).
		Limit(pageSize).Offset(page).Find(&posts).Error
	return posts, err
}

func GetPostsByTagTime(tag string, page int, pageSize int, desc bool) ([]Post, error) {
	posts := make([]Post, 0)
	var err error
	if desc {
		err = global.DB.Model(&Post{}).Preload("User").Where("tag=?", tag).
			Order("created_at DESC").Limit(pageSize).Offset(page).Find(&Post{}).Error
	} else {
		err = global.DB.Model(&Post{}).Preload("User").Where("tag=?", tag).
			Order("created_at").Limit(pageSize).Offset(page).Find(&Post{}).Error
	}
	return posts, err
}

func GetPostsByTagReplies(tag string, page int, pageSize int, desc bool) ([]Post, error) {
	posts := make([]Post, 0)
	var err error

	if desc {
		err = global.DB.Model(&Post{}).Preload("User").Where("tag=?", tag).
			Order("replies DESC").Limit(pageSize).Offset(page).Find(&Post{}).Error
	} else {
		err = global.DB.Model(&Post{}).Preload("User").Where("tag=?", tag).
			Order("replies").Limit(pageSize).Offset(page).Find(&Post{}).Error
	}

	return posts, err
}
