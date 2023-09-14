package router

import (
	"mangosteen/internal/controller"

	"github.com/gin-gonic/gin"
)

func loadControllers() []controller.Controller {
	return []controller.Controller{
		&controller.SessionController{},
		&controller.ValidationCodeController{},
	}
}

func New() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	for _, ctrl := range loadControllers() {
		ctrl.RegisterRoutes(api)
	}

	r.GET("/ping", controller.Ping)

	return r
}
