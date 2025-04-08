package util

import (
	"BBingyan/internal/global"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/smtp"
	"time"
)

func GenerateCode() string {
	return fmt.Sprintf("%4d", rand.Int()%9999)
}

func SendAuthCode(to string, code string) error {
	e := email.NewEmail()
	e.From = "BBingyan <bbingyan@qq.com>"
	e.To = []string{to}
	e.Subject = "Auth Code"
	text := fmt.Sprintf("【BBingyan】验证码: %s ,\n请在5分钟内进行验证,如非本人操作请忽略.\n感谢使用Wumi's BBingyan!", code)
	e.Text = []byte(text)

	return e.Send("smtp.qq.com:25", smtp.PlainAuth("", "bbingyan@qq.com", global.AuthorizationCode, "smtp.qq.com"))
}

func HashPwd(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(b), err
}

func ParsePwd(hashed string, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
}

func GenerateJWT(em string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": em,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	signedToken, err := token.SignedString(global.Key)

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
		return global.Key, nil
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
