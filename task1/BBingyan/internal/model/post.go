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

func GetPostsByEmail(email string) ([]Post, error) {
	posts := make([]Post, 0)
	err := global.DB.Model(&Post{}).Where("author=?", email).Find(&posts).Error
	return posts, err
}

func GetPostById(id int) (*Post, error) {
	post := &Post{}
	err := global.DB.Model(&Post{}).Where("id=?", id).First(post).Error
	return post, err
}

func GetPostLikes(id int) (int, error) {
	var likes int
	err := global.DB.Model(&Post{}).Select("likes").Where("id=?", id).First(&likes).Error
	return likes, err
}
