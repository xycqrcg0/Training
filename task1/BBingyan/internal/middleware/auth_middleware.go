package middleware

import (
	"BBingyan/internal/config"
	"BBingyan/internal/util"
	"github.com/labstack/echo/v4"
)

func CheckJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			path := c.Request().URL.Path
			for _, skipPath := range config.Config.JWT.Skipper {
				if skipPath == path {
					return next(c)
				}
			}

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.ErrBadRequest
			}

			email, err := util.ParseJWT(token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("identification", email)
			return next(c)
		}
	}
}
