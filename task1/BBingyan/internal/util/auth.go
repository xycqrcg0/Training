package util

import (
	"BBingyan/internal/global"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
)

func GenerateCode() string {
	return fmt.Sprintf("%4d", rand.Int()%9999)
}

func SendAuthCode(to string, code string) error {
	e := email.NewEmail()
	e.From = "Wumi's BBingyan <bbingyan@qq.com>"
	e.To = []string{to}
	e.Subject = "Auth Code"
	text := fmt.Sprintf("您的验证码是 %s ,请在1分钟内进行验证.\n感谢使用Wumi's BBingyan!", code)
	e.Text = []byte(text)

	return e.Send("smtp.qq.com:25", smtp.PlainAuth("", "bbingyan@qq.com", global.AuthorizationCode, "smtp.qq.com"))
}
