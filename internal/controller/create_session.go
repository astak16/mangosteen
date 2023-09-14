package controller

import (
	"log"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"mangosteen/sql/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSession(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数有问题"})
		return
	}

	q := database.NewQuery()
	_, err := q.FindValidationCode(c, queries.FindValidationCodeParams{
		Email: requestBody.Email,
		Code:  requestBody.Code,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "验证码无效"})
		return
	}

	user, err := q.FindUserByEmail(c, requestBody.Email)
	if err != nil {
		user, err = q.CreateUser(c, requestBody.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "请稍后再试"})
			log.Println("[CreateSession Fail]", err)
			return
		}
	}

	user_id := int(user.ID)
	jwt, err := jwt_helper.GenerateJWT(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请稍后再试"})
		log.Println("[CreateSession Fail]", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "jwt": jwt, "userId": user_id})
}
