package controller

import (
	"fmt"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MeController struct{}

func (ctrl *MeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.GET("/me", ctrl.Create)
}
func (ctrl *MeController) Create(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 {
		c.JSON(401, gin.H{
			"message": "无效的 JWT",
		})
		return
	}
	jwtString := auth[7:]
	if len(jwtString) == 0 {
		c.JSON(401, gin.H{
			"message": "无效的 JWT",
		})
		return
	}
	token, err := jwt_helper.Parse(jwtString)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "无效的 JWT",
		})
		return
	}
	m, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{
			"message": "无效的 JWT",
		})
		return
	}
	userId := int32(m["user_id"].(float64))

	q := database.NewQuery()
	user, err := q.FindUserById(c, int32(userId))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "无效的 JWT",
		})
		return
	}
	fmt.Println(user.ID)
	c.JSON(200, gin.H{
		"message":  "ok",
		"resource": user,
	})
}
func (ctrl *MeController) Destroy(c *gin.Context)  {}
func (ctrl *MeController) Update(c *gin.Context)   {}
func (ctrl *MeController) Get(c *gin.Context)      {}
func (ctrl *MeController) GetPaged(c *gin.Context) {}
