package router

import (
	"mangosteen/internal/controller"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", controller.Ping)
	r.POST("/api/v1/validation_codes", controller.CreateValidationCode)

	return r
}
