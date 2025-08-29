package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/server/handler"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type authHandlerMocks struct {
	userService *MockuserService
}

func newAuthHandler(t *testing.T) (*gin.Engine, authHandlerMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userService := NewMockuserService(ctrl)
	authHandler := handler.NewAuthHandler(userService)

	engine := gin.New()
	engine.POST("/users", authHandler.RegisterUser)

	mocks := authHandlerMocks{
		userService: userService,
	}

	return engine, mocks
}

func TestAuthHandler_RegisterUser(t *testing.T) {
	registerRequest := request.RegisterRequest{
		BasicAuthRequest: &request.BasicAuthRequest{
			Email:    "oleksandr.khmil@gmail.com",
			Password: "strong-password",
		},
		FullName: "full-name",
	}

	rawRegisterRequest, err := json.Marshal(registerRequest)
	require.NoError(t, err)

	t.Run("It should respond with 400 status if received invalid request", func(t *testing.T) {
		engine, _ := newAuthHandler(t)

		badRegisterRequest := registerRequest
		badRegisterRequest.BasicAuthRequest = &request.BasicAuthRequest{
			Email:    registerRequest.Email,
			Password: "weak",
		}

		rawBadRegisterRequest, err := json.Marshal(badRegisterRequest)
		require.NoError(t, err)

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/users",
			bytes.NewReader(rawBadRegisterRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		response := recorder.Result()
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		expectedResponse := `{
			"code": 400,
			"message": "Invalid Request"
		}`

		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("It should respond with 422 status if received invalid storage operation error", func(t *testing.T) {
		engine, mocks := newAuthHandler(t)

		mocks.userService.
			EXPECT().
			CreateUser(gomock.Any(), registerRequest).
			Return(service.NewErrUserAlreadyExists("msg", "op-name"))

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/users",
			bytes.NewReader(rawRegisterRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		response := recorder.Result()
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)

		expectedResponse := `{
			"code": 422,
			"message": "msg"
		}`

		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("It should create an user", func(t *testing.T) {
		engine, mocks := newAuthHandler(t)

		mocks.userService.
			EXPECT().
			CreateUser(gomock.Any(), registerRequest).
			Return(nil)

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/users",
			bytes.NewReader(rawRegisterRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		response := recorder.Result()
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, `"Successfully registered"`, string(responseBody))
	})
}
