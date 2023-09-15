package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"mangosteen/internal/jwt_helper"
	"mangosteen/sql/queries"
	"mangosteen/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	tc := TagController{}
	tc.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tags", strings.NewReader(`{
		"name": "同情",
		"kind": "expenses",
		"sign": "-"
	}`))

	u := signIn(t, req)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource queries.Tag `json:"resource"`
		Message  string      `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u.ID, j.Resource.UserID)
	assert.Equal(t, "同情", j.Resource.Name)
	assert.Nil(t, j.Resource.DeletedAt)
}

func TestUpdateTag(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	tc := TagController{}
	tc.RegisterRoutes(r.Group("/api"))
	var req *http.Request
	// u := signIn(t, req)
	randNumber, err := utils.RandNumber(10)
	if err != nil {
		t.Fatal(err)
	}
	u, err := q.CreateUser(context.Background(), fmt.Sprintf(`%s@qq.com`, randNumber))
	if err != nil {
		t.Fatal(err)
	}
	jwtString, _ := jwt_helper.GenerateJWT(int(u.ID))

	tag, err := q.CreateTag(context.Background(), queries.CreateTagParams{
		UserID: u.ID,
		Kind:   "expenses",
		Name:   "同情",
		Sign:   "-",
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/api/v1/tags/%d", tag.ID), strings.NewReader(`{
		"name": "吃饭"
	}`))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource queries.Tag `json:"resource"`
		Message  string      `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u.ID, j.Resource.UserID)
	assert.Equal(t, "吃饭", j.Resource.Name)
	assert.Equal(t, "-", j.Resource.Sign)
	assert.Nil(t, j.Resource.DeletedAt)
}
