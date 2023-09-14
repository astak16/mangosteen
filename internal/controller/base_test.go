package controller

import (
	"mangosteen/initialize"
	"mangosteen/internal/database"
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

	q = database.NewQuery()

	return func(t *testing.T) {
		database.Close()
	}
}
