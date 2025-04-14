package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"identification": email,
		"role":           "user",
		"exp":            time.Now().Add(time.Hour * 2).Unix(),
	})

	signedToken, err := token.SignedString([]byte("key"))

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
			return []byte("key"), nil
		})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return jwt.MapClaims{}, err
}
