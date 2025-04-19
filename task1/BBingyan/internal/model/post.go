package model

import (
	"BBingyan/internal/global"
	"gorm.io/gorm"
	"time"
)

//id是能暴露的吗

type Post struct {
	Id        int    `gorm:"id;primaryKey"`
	Author    string `gorm:"author;not null"`
	Title     string `gorm:"title;not null"`
	Content   string `gorm:"content;not null"`
	Likes     int    `gorm:"likes"`
	Replies   int    `gorm:"replies"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func AddPost(newPost *Post) error {
	err := global.DB.Model(&Post{}).Create(newPost).Error
	return err
}

func DeletePostById(id int) error {
	err := global.DB.Model(&Post{}).Where("id=?", id).Delete(&Post{}).Error
	return err
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
