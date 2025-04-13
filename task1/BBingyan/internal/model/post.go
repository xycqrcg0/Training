package model

import (
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
