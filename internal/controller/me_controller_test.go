package controller

import (
	"context"
	"encoding/json"
	"mangosteen/internal/jwt_helper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeController(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	mc := MeController{}
	mc.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestGetMeWithJwt(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	mc := MeController{}
	mc.RegisterRoutes(r.Group("/api"))

	u, err := q.CreateUser(context.Background(), "11dawdafst1@qq.com")
	if err != nil {
		t.Fatal(err)
	}
	jwtString, err := jwt_helper.GenerateJWT(int(u.ID))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource struct {
			ID    int32  `json:"id"`
			Email string `json:"email"`
		} `json:"resource"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "ok", j.Message)
	assert.Equal(t, u.ID, j.Resource.ID)
	assert.Equal(t, u.Email, j.Resource.Email)
}
