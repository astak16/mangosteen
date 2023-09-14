package controller

import (
	"mangosteen/sql/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MeController struct{}

func (ctrl *MeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.GET("/me", ctrl.Get)
}
func (ctrl *MeController) Get(c *gin.Context) {
	me, _ := c.Get("me")

	user, ok := me.(queries.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "无效的 JWT",
		})
		return
	}

	c.JSON(200, gin.H{
		"message":  "ok",
		"resource": user,
	})
}
func (ctrl *MeController) Destroy(c *gin.Context)  {}
func (ctrl *MeController) Update(c *gin.Context)   {}
func (ctrl *MeController) Create(c *gin.Context)   {}
func (ctrl *MeController) GetPaged(c *gin.Context) {}
