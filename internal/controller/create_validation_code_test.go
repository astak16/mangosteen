package controller

import (
	"context"
	"mangosteen/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateValidationCode(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	vcc := ValidationCodeController{}
	vcc.RegisterRoutes(r.Group("/api"))

	email := "1500846601@qq.com"
	viper.Set("viper", &config.ViperConfig{
		Host: "localhost",
		Port: 1025,
	})

	count1, _ := q.CountValidationCodes(context.Background(), email)

	w := httptest.NewRecorder()
	data := strings.NewReader(`{"email":"` + email + `"}`)
	req, _ := http.NewRequest("POST", "/api/v1/validation_codes", data)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	count2, _ := q.CountValidationCodes(context.Background(), email)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, count1+1, count2)
}
