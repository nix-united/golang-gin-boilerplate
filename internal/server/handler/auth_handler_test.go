package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"
	"github.com/nix-united/golang-gin-boilerplate/internal/server/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type authHandlerMocks struct {
	userService     *MockuserService
	passwordService *MockpasswordService
}

func newAuthHandler(t *testing.T) (*gin.Engine, authHandlerMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userService := NewMockuserService(ctrl)
	passwordService := NewMockpasswordService(ctrl)
	authHandler, err := handler.NewAuthHandler(handler.AuthHandlerConfig{
		ApplicationName:         "auth_handler_test",
		JWTSecret:               "auth_handler_secret",
		JWTTokenDuration:        time.Minute,
		JWTTokenRefreshDuration: time.Minute,
		UserService:             userService,
		PasswordService:         passwordService,
	})
	require.NoError(t, err)

	engine := gin.New()
	engine.POST("/register", authHandler.RegisterUser)
	engine.POST("/login", authHandler.LoginUser)
	engine.POST("/refresh", authHandler.Middleware, authHandler.RefreshUserToken)

	mocks := authHandlerMocks{
		userService:     userService,
		passwordService: passwordService,
	}

	return engine, mocks
}

func TestAuthHandler_RegisterUser(t *testing.T) {
	registerRequest := request.RegisterRequest{
		BasicAuthRequest: request.BasicAuthRequest{
			Email:    "name.surname@gmail.com",
			Password: "strong-password",
		},
		FullName: "full-name",
	}

	rawRegisterRequest, err := json.Marshal(registerRequest)
	require.NoError(t, err)

	t.Run("It should respond with 400 status if received invalid request", func(t *testing.T) {
		engine, _ := newAuthHandler(t)

		badRegisterRequest := registerRequest
		badRegisterRequest.BasicAuthRequest = request.BasicAuthRequest{
			Email:    registerRequest.Email,
			Password: "weak",
		}

		rawBadRegisterRequest, err := json.Marshal(badRegisterRequest)
		require.NoError(t, err)

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/register",
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
			"code": "bad_request",
			"message": "Invalid request"
		}`

		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("It should respond with 409 status if received invalid storage operation error", func(t *testing.T) {
		engine, mocks := newAuthHandler(t)

		mocks.userService.
			EXPECT().
			CreateUser(gomock.Any(), registerRequest).
			Return(domain.ErrAlreadyExists)

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/register",
			bytes.NewReader(rawRegisterRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		response := recorder.Result()
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, response.StatusCode)

		expectedResponse := `{
			"code": "already_exists",
			"message": "Such user already exists"
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
			"https://example.com/register",
			bytes.NewReader(rawRegisterRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		httpResponse := recorder.Result()
		defer httpResponse.Body.Close()

		responseBody, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

		wantResponse := response.MessageResponse{
			Message: "Successfully registered",
		}

		var gotResponse response.MessageResponse
		err = json.Unmarshal(responseBody, &gotResponse)
		require.NoError(t, err)

		assert.Equal(t, wantResponse, gotResponse)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	storedUser := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "name.surname@gmail.com",
		Password: "strong-password",
		FullName: "name.surname",
	}

	invalidLoginRequest := request.BasicAuthRequest{
		Email:    storedUser.Email,
		Password: "wrong-password",
	}

	rawInvalidLoginRequest, err := json.Marshal(invalidLoginRequest)
	require.NoError(t, err)

	validLoginRequest := request.BasicAuthRequest{
		Email:    storedUser.Email,
		Password: storedUser.Password,
	}

	rawValidLoginRequest, err := json.Marshal(validLoginRequest)
	require.NoError(t, err)

	t.Run("It should return an error when password is wrong", func(t *testing.T) {
		engine, mocks := newAuthHandler(t)

		mocks.userService.
			EXPECT().
			GetUserByEmail(gomock.Any(), validLoginRequest.Email).
			Return(storedUser, nil)

		mocks.passwordService.
			EXPECT().
			VerifyPassword(storedUser.Password, invalidLoginRequest.Password).
			Return(errors.New("invalid password"))

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/login",
			bytes.NewReader(rawInvalidLoginRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		httpResponse := recorder.Result()
		defer httpResponse.Body.Close()

		responseBody, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)

		var gotResponse response.ErrorResponse
		err = json.Unmarshal(responseBody, &gotResponse)
		require.NoError(t, err)

		wantResponse := response.ErrorResponse{
			Code:    response.CodeAccessDenied,
			Message: "Invalid email or password",
		}

		assert.Equal(t, wantResponse, gotResponse)
	})

	t.Run("It should login an user when password is correct", func(t *testing.T) {
		engine, mocks := newAuthHandler(t)

		mocks.userService.
			EXPECT().
			GetUserByEmail(gomock.Any(), validLoginRequest.Email).
			Return(storedUser, nil)

		mocks.passwordService.
			EXPECT().
			VerifyPassword(storedUser.Password, validLoginRequest.Password).
			Return(nil)

		httpRequest := httptest.NewRequest(
			http.MethodPost,
			"https://example.com/login",
			bytes.NewReader(rawValidLoginRequest),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httpRequest)

		httpResponse := recorder.Result()
		defer httpResponse.Body.Close()

		responseBody, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

		var gotResponse response.AuthTokenResponse
		err = json.Unmarshal(responseBody, &gotResponse)
		require.NoError(t, err)

		assert.NotEmpty(t, gotResponse.AccessToken)
		assert.NotEmpty(t, gotResponse.RefreshToken)
		assert.NotEmpty(t, gotResponse.ExpiresIn)
		assert.Equal(t, "Bearer", gotResponse.TokenType)
	})
}
