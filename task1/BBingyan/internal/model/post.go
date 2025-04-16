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

func LikePost(id int, likes int) error {
	err := global.DB.Model(&Post{}).Where("id=?", id).Update("likes", gorm.Expr("likes+%d", likes)).Error
	return err
}
