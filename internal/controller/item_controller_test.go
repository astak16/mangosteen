package controller

import (
	"context"
	"encoding/json"
	"mangosteen/api"
	"mangosteen/sql/queries"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/nav-inc/datetime"
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

func TestGetBalance(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	if err := q.DeleteAllItems(context.Background()); err != nil {
		t.Fatal(err)
	}

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/items/balance", nil)

	u := signIn(t, req)

	for i := 0; i < 10; i++ {
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
	var j api.GetBalanceResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int32(100000), j.Expenses)
	assert.Equal(t, int32(0), j.Income)
	assert.Equal(t, int32(-100000), j.Balance)
}

func TestGetBalanceWithTime(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	if err := q.DeleteAllItems(context.Background()); err != nil {
		t.Fatal(err)
	}

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/items/balance?happened_after="+url.QueryEscape("2020-01-01T00:00:00+0800")+"&happened_before="+url.QueryEscape("2020-01-02T00:00:00+0800"), nil)

	u := signIn(t, req)

	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2019-01-01T00:00:00+0800", time.Local)
		if _, err := q.CreateItem(context.Background(), queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2020-01-01T12:00:00+0800", time.Local)
		if _, err := q.CreateItem(context.Background(), queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2020-01-10T12:00:00:0800", time.Local)
		if _, err := q.CreateItem(context.Background(), queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Fatal(err)
		}
	}

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var j api.GetBalanceResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int32(10000*3), j.Expenses)
	assert.Equal(t, int32(0), j.Income)
	assert.Equal(t, int32(-10000*3), j.Balance)
}
