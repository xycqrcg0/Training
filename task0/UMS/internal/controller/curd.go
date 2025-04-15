package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ums/internal/controller/params"
	"ums/internal/models"
	"ums/internal/utils"
)

func GetInfo(c echo.Context) error {
	email := c.Get("identification").(string)
	role := c.Get("role").(string)
	if role != "user" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not a user",
		})
	}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	res := &params.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "",
		Data:   res,
	})
}

// UpdateUser 邮箱不让修改
func UpdateUser(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "user" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not a user",
		})
	}
	email := c.Get("identification").(string)

	data := params.UserUpdateRequest{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	//检查一下password格式

	user, err := models.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Password != "" {
		hashed, err := utils.HashPassword(data.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &params.Response{
				Status: false,
				Msg:    err.Error(),
			})
		}
		user.Password = hashed
	}

	if err := models.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "Update successfully",
	})
}

func DeleteUser(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "user" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not a user",
		})
	}
	email := c.Get("identification").(string)

	if err := models.DeleteUser(email); err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "Delete successfully",
	})
}
