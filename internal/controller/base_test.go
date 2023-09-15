package controller

import (
	"context"
	"fmt"
	"mangosteen/config"
	"mangosteen/initialize"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"mangosteen/internal/middleware"
	"mangosteen/sql/queries"
	"mangosteen/utils"
	"net/http"
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

func signIn(t *testing.T, req *http.Request) queries.User {
	randNumber, err := utils.RandNumber(10)
	if err != nil {
		t.Fatal(err)
	}
	u, err := q.CreateUser(context.Background(), fmt.Sprintf(`%s@qq.com`, randNumber))
	if err != nil {
		t.Fatal(err)
	}
	jwtString, _ := jwt_helper.GenerateJWT(int(u.ID))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	return u
}
