package email

import (
	"fmt"
	"log"
	"mangosteen/global"

	"gopkg.in/gomail.v2"
)

func Send() {
	m := newMessage("1500846601@qq.com", "Hello", "Hello, <b>你好</b>!")
	d := newDialer()
	if err := d.DialAndSend(m); err != nil {
		log.Fatalln(err)
	}
}

func newDialer() *gomail.Dialer {
	viperConfig := global.ViperConfig
	d := gomail.NewDialer(viperConfig.Host, viperConfig.Port, viperConfig.Username, viperConfig.Password)
	return d
}

func SendValidationCode(email, code string) error {
	m := newMessage(email, fmt.Sprintf("[%s] 验证码", code), fmt.Sprintf("您的验证码为: %s", code))
	d := newDialer()
	return d.DialAndSend(m)
}

func newMessage(to, subject, body string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", "1500846601@qq.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return m
}
