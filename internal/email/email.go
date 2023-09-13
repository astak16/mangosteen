package email

import (
	"fmt"
	"log"
	"mangosteen/global"

	"gopkg.in/gomail.v2"
)

func Send() {

	fmt.Println(global.ViperConfig, "uccs")

	m := gomail.NewMessage()
	m.SetHeader("From", "1500846601@qq.com")
	m.SetHeader("To", "1500846601@qq.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello 张凯君 ucscs")

	viperConfig := global.ViperConfig
	d := gomail.NewDialer(viperConfig.Host, viperConfig.Port, viperConfig.Username, viperConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Fatalln(err)
	}
}
