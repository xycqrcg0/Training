package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        int    `gorm:"id;PrimaryKey"`
	Email     string `gorm:"email;not null"`
	Name      string `gorm:"name;not null"`
	Password  string `gorm:"password;not null"`
	Signature string `gorm:"signature;default:你好，世界"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
