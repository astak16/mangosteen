package router

import (
	"mangosteen/config"
	"mangosteen/internal/controller"
	"mangosteen/internal/middleware"

	"github.com/gin-gonic/gin"
)

func loadControllers() []controller.Controller {
	return []controller.Controller{
		&controller.PingController{},
		&controller.MeController{},
		&controller.TagController{},
		&controller.ItemController{},
		&controller.SessionController{},
		&controller.ValidationCodeController{},
	}
}

func New() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Me(config.WhiteList))

	api := r.Group("/api")

	for _, ctrl := range loadControllers() {
		ctrl.RegisterRoutes(api)
	}

	return r
}
