package controller

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ums/internal/controller/params"
	"ums/internal/models"
)

func GetInfo(c echo.Context) error {
	email := c.Get("email").(string)

	user, err := models.GetUserByEmail(email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	res := &params.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateUser 邮箱不让修改
func UpdateUser(c echo.Context) error {
	email := c.Get("email").(string)
	name := c.QueryParam("name")
	password := c.QueryParam("pwd")
	//检查一下password格式

	user, err := models.GetUserByEmail(email)
	if err != nil {
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

	if err := models.UpdateUser(user); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func DeleteUser(c echo.Context) error {
	email := c.Get("email").(string)

	if err := models.DeleteUser(email); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}
