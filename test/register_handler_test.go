package test

import (
	"basic_server/test/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestRegisterAttemptWithEmptyRequestPayload(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", nil)
	req.Header.Add("Content-Type", "application/json")

	service.TestEngine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegisterAttemptWithEmptyEmailField(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "",
		"password": "test password",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestEngine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestSuccessfulRegisterAttempt(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"email":     "test@test.com",
		"password":  "test password",
		"full_name": "test test",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestEngine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessfulRegisterAttemptWithoutFullNameField(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "test@test.com",
		"password": "test password",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestEngine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
