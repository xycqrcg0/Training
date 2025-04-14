package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ums/internal/controller/params"
)

func AdminLogin(c echo.Context) error {
	admin := &params.UserLoginRequest{}
	if err := c.Bind(admin); err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

}
