package controller

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ums/internal/global"
	"ums/internal/models"
)

func GetInfo(c echo.Context) error {
	email := c.Get("email").(string)

	user := &models.User{}
	if err := global.DB.Model(&user).Where("email=?", email).First(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	res := &models.ResUser{
		Name:  user.Name,
		Email: user.Email,
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateUser 邮箱不让修改
func UpdateUser(c echo.Context) error {
	email := c.Get("email").(string)
	name := c.QueryParam("name")
	password := c.QueryParam("pwd")
	//检查一下password格式

	user := &models.User{}
	if err := global.DB.Model(&user).Where("email=?", email).First(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	if name != "" {
		user.Name = name
	}
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return echo.ErrInternalServerError
		}
		user.Password = string(hashed)
	}

	if err := global.DB.Model(&user).Where("email=?", email).Updates(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func DeleteUser(c echo.Context) error {
	email := c.Get("email").(string)

	if err := global.DB.Model(&models.User{}).Where("email=?", email).Delete(&models.User{}).Error; err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}
