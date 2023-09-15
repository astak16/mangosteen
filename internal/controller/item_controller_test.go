package controller

import (
	"context"
	"encoding/json"
	"mangosteen/api"
	"mangosteen/sql/queries"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

	u := signIn(t, req)
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
	signIn(t, req)
	r.ServeHTTP(w, req)

	assert.Equal(t, 422, w.Code)
}

func TestGetPagesItems(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	if err := q.DeleteAllItems(context.Background()); err != nil {
		t.Fatal(err)
	}

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/items", nil)

	u := signIn(t, req)

	for i := 0; i < int(ic.PerPage-2); i++ {
		if _, err := q.CreateItem(context.Background(), queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: time.Now(),
		}); err != nil {
			t.Fatal(err)
		}
	}

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var j api.GetPagedItemsResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ic.PerPage-2, int32(len(j.Resources)))
	// assert.Equal(t, int32(100), j.Resource.Amount)
}
