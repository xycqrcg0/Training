package controller

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ums/internal/global"
	"ums/internal/models"
	"ums/internal/utils"
)

// Login 要求email不重复，用email登录
func Login(c echo.Context) error {
	userInfo := &models.LReqUser{}
	if err := c.Bind(&userInfo); err != nil {
		return echo.ErrBadRequest
	}

	user := &models.User{}
	if err := global.DB.Model(&models.User{}).Where("email=?", userInfo.Email).First(&user).Error; err != nil {
		return echo.ErrForbidden
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
	userInfo := &models.RReqUser{}
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
	var count int64
	if err := global.DB.Model(&models.User{}).Where("email=?", userInfo.Email).Count(&count).Error; err != nil {
		return echo.ErrInternalServerError
	}
	if count > 0 {
		return echo.ErrBadRequest
	}

	//填入数据库
	newUser := &models.User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
	}
	if err := global.DB.Create(&newUser).Error; err != nil {
		return echo.ErrInternalServerError
	}

	//返回一个token
	token, err := utils.GenerateJWT(newUser.Email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, &models.Message{Message: token})
}
