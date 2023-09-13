package controller_test

import (
	"context"
	"mangosteen/config"
	"mangosteen/initialize"
	"mangosteen/internal/database"
	"mangosteen/internal/router"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateValidationCode(t *testing.T) {
	r := router.New()
	email := "1500846601@qq.com"
	viper.Set("viper", &config.ViperConfig{
		Host: "localhost",
		Port: 1025,
	})
	database.Connect()
	defer database.Close()

	q := database.NewQuery()
	count1, _ := q.CountValidationCodes(context.Background(), email)

	initialize.InitViper()

	w := httptest.NewRecorder()
	data := strings.NewReader(`{"email":"` + email + `"}`)
	req, _ := http.NewRequest("POST", "/api/v1/validation_codes", data)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	count2, _ := q.CountValidationCodes(context.Background(), email)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, count1+1, count2)
}
