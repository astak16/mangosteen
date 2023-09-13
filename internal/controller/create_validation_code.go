package controller

import (
	"log"
	"mangosteen/internal/database"
	"mangosteen/internal/email"
	"mangosteen/sql/queries"
	"mangosteen/utils"

	"github.com/gin-gonic/gin"
)

func CreateValidationCode(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(400, "参数错误")
		return
	}

	str, err := utils.RandNumber(4)
	if err != nil {
		log.Println("[CreateValidationCode 创建验证码失败]", err)
		c.String(400, "发送失败")
		return
	}

	q := database.NewQuery()
	_, err = q.CreateValidationCode(c, queries.CreateValidationCodeParams{
		Email: body.Email,
		Code:  str,
	})
	if err != nil {
		log.Println("[CreateValidationCode Fail]", err)
		c.String(400, "创建失败")
		return
	}

	err = email.SendValidationCode(body.Email, str)
	if err != nil {
		log.Println("[SendValidationCode Fail]", err)
		c.String(500, "发送失败")
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
