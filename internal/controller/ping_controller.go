package controller

import (
	"github.com/gin-gonic/gin"
)

type PingController struct{}

func (ctrl *PingController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.GET("/ping", ctrl.Create)
}
func (ctrl *PingController) Create(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
func (ctrl *PingController) Destroy(c *gin.Context)  {}
func (ctrl *PingController) Update(c *gin.Context)   {}
func (ctrl *PingController) Get(c *gin.Context)      {}
func (ctrl *PingController) GetPaged(c *gin.Context) {}
