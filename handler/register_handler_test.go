package handler

import (
	"basic_server/internal/model"
	"basic_server/repository/mocks"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"

	"basic_server/service"

	"github.com/gin-gonic/gin"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := model.User{
		Email:    "test@test.com",
		Password: "11111111",
		FullName: "test name",
	}
	userRepoMock := mocks.NewUserRepositoryMock(&testData)

	server := gin.New()
	server.POST(
		"/users",
		NewAuthHandler(service.NewUserService(userRepoMock)).RegisterUser,
	)

	recorder := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewBuffer([]byte(`{"email":"test@test.com","password":"test"}`)),
	)
	req.Header.Add("Content-Type", "application/json")

	server.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
