package controller

import (
	"BBingyan/internal/global"
	"BBingyan/internal/model"
	"BBingyan/internal/util"
	"errors"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Ruser struct {
	Code     string `json:"code"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Uuser struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Signature string `json:"signature"`
}

func RegisterForCode(c echo.Context) error {
	//在query给个email就行
	//要看email是否符合格式
	email := c.Request().URL.Query().Get("email")

	redisKey := "register:email:" + email

	//先查该email是否在redis里有记录
	codeKey := "code:" + email
	_, err0 := global.RedisDB.Get(codeKey).Result()
	if err0 == nil {
		return c.JSON(http.StatusForbidden, "")
	} else {
		if !errors.Is(err0, redis.Nil) {
			global.Errors.Error("Fail to read code from redis")
			return c.JSON(http.StatusInternalServerError, "")
		}
	}

	//查email是否已存在
	var count int64
	if err := global.DB.Model(&model.User{}).Where("email=?", email).Count(&count).Error; err != nil {
		global.Errors.WithField("target", "email").Error("Fail to search postgres")
		return c.JSON(http.StatusInternalServerError, "")
	}
	if count != 0 {
		global.Infos.WithField("email", email).Info("repeating email")
		return c.JSON(http.StatusForbidden, "")
	}

	//生成验证码并发送
	//存下
	code := util.GenerateCode()
	_, err := global.RedisDB.TxPipelined(func(pipe redis.Pipeliner) error {
		//限制发送
		pipe.Set(codeKey, "", time.Minute*1)
		//记录code
		pipe.Set(redisKey, code, time.Minute*5)
		return nil
	})
	if err != nil {
		global.Errors.Error("Fail to write in redis")
		return c.JSON(http.StatusInternalServerError, "")
	}

	if err := util.SendAuthCode(email, code); err != nil {
		global.Errors.WithField("error", err).Warn("Fail to send email")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, "")
}

func Register(c echo.Context) error {
	data := &Ruser{}

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	//数据都先检查一遍
	if data.Email == "" || data.Code == "" || data.Password == "" || data.Name == "" {
		return c.JSON(http.StatusBadRequest, "")
	}

	//检查code
	redisKey := "register:email:" + data.Email
	com, err := global.RedisDB.Get(redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return c.JSON(http.StatusBadRequest, "")
		}
		global.Errors.Error("Fail to read from redis for code")
		//邮箱对应的code无/过期
		return c.JSON(http.StatusInternalServerError, "")
	}
	if data.Code != com {
		return c.JSON(http.StatusBadRequest, "")
	}

	//能到这说明验证码过了，可以注册
	data.Password, err = util.HashPwd(data.Password)
	if err != nil {
		global.Errors.Warn("Fail to hash password")
		return c.JSON(http.StatusInternalServerError, "")
	}

	newUser := &model.User{
		Email:    data.Email,
		Name:     data.Name,
		Password: data.Password,
	}

	if err := global.DB.Create(&newUser).Error; err != nil {
		global.Errors.Error("Fail to insert a user into postgres")
		return c.JSON(http.StatusInternalServerError, "")
	}

	token, err := util.GenerateJWT(newUser.Email)
	if err != nil {
		global.Errors.Warn("Fail to generate jwt")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusCreated, token)
}

func LoginForCode(c echo.Context) error {
	//在query给个email就行
	//要看email是否符合格式
	email := c.Request().URL.Query().Get("email")

	redisKey := "login:email:" + email

	//先查该email是否在redis里有记录
	codeKey := "code:" + email
	_, err0 := global.RedisDB.Get(codeKey).Result()
	if err0 == nil {
		return c.JSON(http.StatusForbidden, "")
	} else {
		if !errors.Is(err0, redis.Nil) {
			global.Errors.Error("Fail to read code from redis")
			return c.JSON(http.StatusInternalServerError, "")
		}
	}
	//生成验证码并发送
	//存下
	code := util.GenerateCode()
	_, err := global.RedisDB.TxPipelined(func(pipe redis.Pipeliner) error {
		//限制发送
		pipe.Set(codeKey, "", time.Minute*1)
		//记录code
		pipe.Set(redisKey, code, time.Minute*5)
		return nil
	})
	if err != nil {
		global.Errors.Error("Fail to write in redis")
		return c.JSON(http.StatusInternalServerError, "")
	}

	if err := util.SendAuthCode(email, code); err != nil {
		global.Errors.WithField("error", err).Warn("Fail to send email")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, "")
}

func Login(c echo.Context) error {
	style := c.Param("style")
	if style != "v1" && style != "v2" {
		return c.JSON(http.StatusNotFound, "")
	}
	//v1->code;v2->password

	data := &Ruser{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	if style == "v1" {
		//code
		redisKey := "login:email:" + data.Email
		com, err := global.RedisDB.Get(redisKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return c.JSON(http.StatusBadRequest, "")
			}
			global.Errors.Error("Fail to read from redis for code")
			//邮箱对应的code无/过期
			return c.JSON(http.StatusInternalServerError, "")
		}
		if data.Code != com {
			return c.JSON(http.StatusBadRequest, "")
		}
	} else {
		//password
		user := &model.User{}
		if err := global.DB.Model(&model.User{}).Where("email=?", data.Email).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusForbidden, "")
			}
			global.Errors.Error("Fail to read user from postgres")
			return c.JSON(http.StatusInternalServerError, "")
		}
		if err := util.ParsePwd(user.Password, data.Password); err != nil {
			return c.JSON(http.StatusForbidden, "")
		}
	}

	token, err := util.GenerateJWT(data.Email)
	if err != nil {
		global.Errors.Warn("Fail to generate jwt")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, token)
}

func UpdateInfo(c echo.Context) error {
	email := c.Get("email").(string)

	data := &Uuser{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	if data.Name == "" && data.Password == "" && data.Signature == "" {
		return c.JSON(http.StatusBadRequest, "")
	}
	user := &model.User{}
	if err := global.DB.Model(&model.User{}).Where("email=?", email).First(&user).Error; err != nil {
		global.Errors.Error("Fail to read user from postgres")
		return c.JSON(http.StatusInternalServerError, "")
	}
	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Signature != "" {
		user.Signature = data.Signature
	}
	if data.Password != "" {
		b, err := util.GenerateJWT(data.Password)
		if err != nil {
			global.Errors.Error("Fail to hash password")
			return c.JSON(http.StatusInternalServerError, "")
		}
		user.Password = b
	}
	if err := global.DB.Model(&user).Where("email=?", email).Updates(&user).Error; err != nil {
		global.Errors.Error("Fail to update user in postgres")
		return c.JSON(http.StatusInternalServerError, "")
	}
	return c.JSON(http.StatusOK, "")
}
