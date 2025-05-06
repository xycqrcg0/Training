package util

import (
	"BBingyan/internal/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(em string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": em,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	signedToken, err := token.SignedString(config.Config.JWT.Key)

	return "Bearer " + signedToken, err
}

func ParseJWT(signedToken string) (string, error) {
	if len(signedToken) > 7 && signedToken[:7] == "Bearer " {
		signedToken = signedToken[7:]
	} else {
		return "", errors.New("invalid")
	}

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid")
		}
		return config.Config.JWT.Key, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		em, _ := claims["email"].(string)
		return em, nil
	}

	return "", errors.New("expired")
}
