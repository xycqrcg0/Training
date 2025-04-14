package middlewares

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"ums/internal/config"
	"ums/internal/controller/params"
	"ums/internal/utils"
)

func CheckJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			path := c.Request().URL.Path
			log.Println("path:", path)
			for _, skipPath := range config.Config.JWT.Skipper {
				if skipPath == path {
					return next(c)
				}
			}

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusBadRequest, &params.Response{
					Status: false,
					Msg:    "Invalid token",
				})
			}

			claims, err := utils.ParseJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &params.Response{
					Status: false,
					Msg:    "Invalid token",
				})
			}

			c.Set("identification", claims["identification"].(string))
			c.Set("role", claims["role"].(string))

			return next(c)
		}
	}
}
