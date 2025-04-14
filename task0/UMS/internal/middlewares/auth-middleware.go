package middlewares

import (
	"github.com/labstack/echo/v4"
	"ums/internal/utils"
)

func CheckJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.ErrBadRequest
			}

			claims, err := utils.ParseJWT(token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("identification", claims["identification"].(string))
			c.Set("role", claims["role"].(string))
			return next(c)
		}
	}
}
