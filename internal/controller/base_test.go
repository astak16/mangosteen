package controller

import (
	"mangosteen/config"
	"mangosteen/initialize"
	"mangosteen/internal/database"
	"mangosteen/internal/middleware"
	"mangosteen/sql/queries"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
	q *queries.Queries
)

func setupTest(t *testing.T) func(t *testing.T) {
	database.Connect()
	initialize.InitViper()
	r = gin.Default()
	r.Use(middleware.Me(config.WhiteList))

	q = database.NewQuery()

	return func(t *testing.T) {
		database.Close()
	}
}
