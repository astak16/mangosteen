package controller

import (
	"context"
	"encoding/json"
	"mangosteen/internal/jwt_helper"
	"mangosteen/sql/queries"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/items", strings.NewReader(`{
		"amount": 100,
		"kind": "expenses",
		"happened_at": "2020-01-01T00:00:00+08:00",
		"tag_ids": [1, 2]
	}`))

	u, _ := q.CreateUser(context.Background(), "24rd23rr@test.com")
	jwtString, _ := jwt_helper.GenerateJWT(int(u.ID))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource queries.Item `json:"resource"`
		Message  string       `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u.ID, j.Resource.ID)
	assert.Equal(t, int32(100), j.Resource.Amount)
}

func TestCreateItemWithError(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/items", strings.NewReader(`{
		"kind": "expenses",
		"happened_at": "2020-01-01T00:00:00+08:00",
		"tag_ids": [1, 2]
	}`))

	u, _ := q.CreateUser(context.Background(), "24rsdrrrd23rr@test.com")
	jwtString, _ := jwt_helper.GenerateJWT(int(u.ID))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 422, w.Code)
}
