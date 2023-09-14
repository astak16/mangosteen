package controller

import (
	"time"

	"github.com/gin-gonic/gin"
)

type ItemController struct{}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/items", ctrl.Create)
}
func (ctrl *ItemController) Create(c *gin.Context) {
	var body struct {
		Amount     int32     `json:"amount" binding:"required"`
		Kind       string    `json:"kind" binding:"required"`
		HappenedAt time.Time `json:"happened_at" binding:"required"`
		TagIds     []int32   `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(422, gin.H{"error": "请求参数有误"})
		return
	}
}
func (ctrl *ItemController) Destroy(c *gin.Context)  {}
func (ctrl *ItemController) Update(c *gin.Context)   {}
func (ctrl *ItemController) Get(c *gin.Context)      {}
func (ctrl *ItemController) GetPaged(c *gin.Context) {}
