package test

import (
	"basic_server/server/model"
	"basic_server/server/provider"
	"basic_server/test/service"
	"basic_server/test/service/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSuccessfulLogin(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	testPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	service.TestServer().DatabaseDriver().Create(
		&model.User{
			Email:    "test1@test.com",
			Password: string(testPassword),
		},
	)

	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "test1@test.com",
		"password": "test",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessfulRefreshTokenAttempt(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	testPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	user := model.User{
		Email:    "test1@test.com",
		Password: string(testPassword),
	}

	service.TestServer().DatabaseDriver().Create(&user)

	jwtAuth := provider.NewJwtAuth(service.TestServer().DatabaseDriver())

	token, _, _ := jwtAuth.Middleware().TokenGenerator(user)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/refresh", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginAttemptWithInvalidCreadentials(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	testPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	service.TestServer().DatabaseDriver().Create(
		&model.User{
			Email:    "test1@test.com",
			Password: string(testPassword),
		},
	)

	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "test2@test.com",
		"password": "test",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginAttemptWithInvalidEmailFormat(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "test",
		"password": "test",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginAttemptWithEmptyEmailField(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"password": "test",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginAttemptWithEmptyPasswordField(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"email": "test2@test.com",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshTokenAttemptWithInvalidJwtToken(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	testPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	service.TestServer().DatabaseDriver().Create(
		&model.User{
			Email:    "test1@test.com",
			Password: string(testPassword),
		},
	)

	jwtAuth := provider.NewJwtAuth(service.TestServer().DatabaseDriver())

	token, _, _ := jwtAuth.Middleware().TokenGenerator(
		model.User{
			Email:    "test2@test.com",
			Password: "test",
		},
	)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/refresh", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRefreshTokenAttemptWithoutJwtToken(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/refresh", nil)
	req.Header.Add("Accept", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
