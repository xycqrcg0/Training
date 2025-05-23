package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"ums/internal/config"
)

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"identification": email,
		"role":           "user",
		"exp":            time.Now().Add(time.Minute * time.Duration(config.Config.JWT.Exp)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(config.Config.JWT.Key))

	return "Bearer " + signedToken, err
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	if len(tokenString) > 7 || tokenString[:7] == "Bearer" {
		tokenString = tokenString[7:]
	} else {
		return jwt.MapClaims{}, errors.New("not a Bearer token")
	}

	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", errors.New("invalid token")
			}
			return []byte(config.Config.JWT.Key), nil
		})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return jwt.MapClaims{}, err
}
