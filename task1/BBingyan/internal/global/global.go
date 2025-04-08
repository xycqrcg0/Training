package global

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

var (
	AuthorizationCode string

	Key []byte

	DB *gorm.DB

	Errors *logrus.Logger
	Infos  *logrus.Logger

	Logfile *os.File

	RedisDB *redis.Client
)
