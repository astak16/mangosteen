package controller

import (
	"log"
	"mangosteen/internal/database"
	"mangosteen/internal/email"
	"mangosteen/sql/queries"
	"mangosteen/utils"

	"github.com/gin-gonic/gin"
)

type ValidationCodeController struct{}

func (ctrl *ValidationCodeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/validation_codes", ctrl.Create)
}
func (ctrl *ValidationCodeController) Create(c *gin.Context) {
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
		c.String(401, "创建失败")
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
func (ctrl *ValidationCodeController) Destroy(c *gin.Context)  {}
func (ctrl *ValidationCodeController) Update(c *gin.Context)   {}
func (ctrl *ValidationCodeController) Get(c *gin.Context)      {}
func (ctrl *ValidationCodeController) GetPaged(c *gin.Context) {}
