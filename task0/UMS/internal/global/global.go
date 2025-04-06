package global

import (
	"gorm.io/gorm"
	"ums/internal/models"
)

var Configs models.Config

var DB *gorm.DB
