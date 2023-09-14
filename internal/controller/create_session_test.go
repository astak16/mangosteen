package controller

import (
	"context"
	"encoding/json"
	"log"
	"mangosteen/sql/queries"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	sc := SessionController{}
	sc.RegisterRoutes(r.Group("/api"))

	if err := q.DeleteAllUsers(context.Background()); err != nil {
		t.Fatal(err)
	}

	j := queries.CreateValidationCodeParams{
		Email: "15006601@qq.com",
		Code:  "2345",
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
