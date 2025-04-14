package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"ums/internal/controller/params"
	"ums/internal/models"
	"ums/internal/utils"
)

// Login 要求email不重复，用email登录
func Login(c echo.Context) error {
	userInfo := &params.UserLoginRequest{}
	if err := c.Bind(&userInfo); err != nil {
		return echo.ErrBadRequest
	}

	user, err := models.GetUserByEmail(userInfo.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrForbidden
		}
		return echo.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		return echo.ErrForbidden
	}

	//身份认证成功，生成jwt
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &models.Message{Message: token})
}

func Register(c echo.Context) error {
	userInfo := &params.UserRegisterRequest{}
	if err := c.Bind(&userInfo); err != nil {
		return echo.ErrBadRequest
	}

	if userInfo.Password == "" {
		return echo.ErrBadRequest
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	if err != nil {
		return echo.ErrInternalServerError
	}
	userInfo.Password = string(hashed)

	if userInfo.Email == "" || userInfo.Name == "" {
		return echo.ErrBadRequest
	}
	//检查邮箱格式

	//检查邮箱是否存在
	_, err = models.GetUserByEmail(userInfo.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrInternalServerError
		}
	} else {
		return echo.ErrBadRequest
	}

	//填入数据库
	newUser := &models.User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
	}
	if err := models.AddUser(newUser); err != nil {
		return echo.ErrInternalServerError
	}

	//返回一个token
	token, err := utils.GenerateJWT(newUser.Email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, &models.Message{Message: token})
}
