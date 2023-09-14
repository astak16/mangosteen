package controller_test

import (
	"context"
	"encoding/json"
	"log"
	"mangosteen/initialize"
	"mangosteen/internal/database"
	"mangosteen/internal/router"
	"mangosteen/sql/queries"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	r *gin.Engine
	q *queries.Queries
)

func setupTest(t *testing.T) func(t *testing.T) {
	database.Connect()
	r = router.New()
	q = database.NewQuery()

	if err := q.DeleteAllUsers(context.Background()); err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		database.Close()
	}
}

func TestCreateSession(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	initialize.InitViper()
	j := queries.CreateValidationCodeParams{
		Email: "1500846601@qq.com",
		Code:  "11111",
	}

	w := httptest.NewRecorder()

	_, err := q.CreateValidationCode(context.Background(), j)
	if err != nil {
		log.Fatalln(err)
	}
	b, _ := json.Marshal(j)

	user, err := q.CreateUser(context.Background(), j.Email)
	if err != nil {
		log.Fatalln(err)
	}

	req, _ := http.NewRequest("POST", "/api/v1/session", strings.NewReader(string(b)))
	r.ServeHTTP(w, req)

	var responseBody struct {
		Jwt    string `json:"jwt"`
		UserID int32  `json:"userId"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		log.Fatalln(err)
		return
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, user.ID, responseBody.UserID)
}
