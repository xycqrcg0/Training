package global

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	AuthorizationCode string

	Key []byte

	DB *gorm.DB

	Errors *logrus.Logger

	RedisDB *redis.Client

	ErrPostNone          = errors.New("none")
	ErrFollowExisted     = errors.New("has followed")
	ErrFollowNonexistent = errors.New("hasn't followed")
)
