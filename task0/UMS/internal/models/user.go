package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type LReqUser struct {
	Email    string
	Password string
}

type RReqUser struct {
	Name     string
	Email    string
	Password string
}

type ResUser struct {
	Name  string
	Email string
}
