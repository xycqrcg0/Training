package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	AuthorizationCode string

	DB *gorm.DB

	RedisDB *redis.Client
)
