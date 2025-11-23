package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	ginjwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=$GOFILE -destination=auth_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

const identityKey = "id"

type userService interface {
	CreateUser(ctx context.Context, registerRequest request.RegisterRequest) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type passwordService interface {
	VerifyPassword(actual, received string) error
}

type AuthHandlerConfig struct {
	ApplicationName         string
	JWTSecret               string
	JWTTokenDuration        time.Duration
	JWTTokenRefreshDuration time.Duration
	UserService             userService
	PasswordService         passwordService
}

type AuthHandler struct {
	ginJWT          ginjwt.GinJWTMiddleware
	userService     userService
	passwordService passwordService
}

func NewAuthHandler(config AuthHandlerConfig) (*AuthHandler, error) {
	authHandler := &AuthHandler{
		userService:     config.UserService,
		passwordService: config.PasswordService,
	}

	authHandler.ginJWT = ginjwt.GinJWTMiddleware{
		Realm:                 config.ApplicationName,
		Key:                   []byte(config.JWTSecret),
		Timeout:               config.JWTTokenDuration,
		MaxRefresh:            config.JWTTokenRefreshDuration,
		Authenticator:         authHandler.authenticate,
		PayloadFunc:           authHandler.payload,
		HTTPStatusMessageFunc: authHandler.mapHTTPStatusMessage,
		Unauthorized:          authHandler.respondWithUnauthorized,
		LoginResponse:         authHandler.respondWithAuthToken,
		RefreshResponse:       authHandler.respondWithAuthToken,
	}

	if err := authHandler.ginJWT.MiddlewareInit(); err != nil {
		return nil, fmt.Errorf("init gin jwt middleware: %w", err)
	}

	return authHandler, nil
}

// RegisterUser godoc
// @Summary Register
// @Description New user registration
// @ID user-register
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body request.RegisterRequest true "User's email, password, full name"
// @Success 200 {string} string "Successfully registered"
// @Failure 409 {object} response.ErrorResponse
// @Router /users [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var registerRequest request.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.Error(fmt.Errorf("bind: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	if err := registerRequest.Validate(); err != nil {
		c.Error(fmt.Errorf("validate: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	if err := h.userService.CreateUser(c.Request.Context(), registerRequest); err != nil {
		c.Error(fmt.Errorf("create user: %w", err))
		if errors.Is(err, domain.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, response.NewErrorResponse(
				response.CodeAlreadyExists,
				"Such user already exists",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeBadRequest,
			"Oops, something went wrong...",
		))
		return
	}

	c.JSON(http.StatusOK, response.NewMessageResponse("Successfully registered"))
}

// authenticate godoc
// @Summary Authenticate a user
// @Description Perform user login
// @ID user-login
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body request.BasicAuthRequest true "User's credentials"
// @Failure 401 {object} response.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	h.ginJWT.LoginHandler(c)
}

// refresh godoc
// @Summary Refresh token
// @Description Refresh user's token
// @ID refresh-token
// @Tags User Actions
// @Produce json
// @Failure 401 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	h.ginJWT.RefreshHandler(c)
}

func (h *AuthHandler) Middleware(c *gin.Context) {
	h.ginJWT.MiddlewareFunc()(c)
}

// Returned data will be propagated to [AuthHandler.payload] func by lib.
func (h *AuthHandler) authenticate(c *gin.Context) (any, error) {
	var authRequest request.BasicAuthRequest
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		err = errors.Join(fmt.Errorf("bind: %w", err), ginjwt.ErrMissingLoginValues)
		c.Error(err)
		return nil, err
	}

	storedUser, err := h.userService.GetUserByEmail(c.Request.Context(), authRequest.Email)
	if err != nil {
		err = fmt.Errorf("get user by email: %w", err)
		if errors.Is(err, domain.ErrNotFound) {
			err = errors.Join(err, ginjwt.ErrFailedAuthentication)
		}
		c.Error(err)
		return nil, err
	}

	if err := h.passwordService.VerifyPassword(storedUser.Password, authRequest.Password); err != nil {
		err = errors.Join(fmt.Errorf("verify password: %w", err), ginjwt.ErrFailedAuthentication)
		c.Error(err)
		return nil, err
	}

	return storedUser, nil
}

// payload receives data from [AuthHandler.authenticate].
func (h *AuthHandler) payload(data any) jwt.MapClaims {
	user, ok := data.(*model.User)
	if !ok {
		return jwt.MapClaims{}
	}
	return jwt.MapClaims{identityKey: user.ID}
}

// mapHTTPStatusMessage propagates message to [AuthHandler.unauthorized] by lib.
func (h *AuthHandler) mapHTTPStatusMessage(_ *gin.Context, err error) string {
	switch {
	case errors.Is(err, ginjwt.ErrMissingLoginValues):
		return "Email and password are required"
	case errors.Is(err, ginjwt.ErrFailedAuthentication):
		return "Invalid email or password"
	default:
		return "Internal server error"
	}
}

// respondWithUnauthorized receives message from [AuthHandler.httpStatusMessage].
func (h *AuthHandler) respondWithUnauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, response.NewErrorResponse(response.CodeAccessDenied, message))
}

// respondWithAuthToken maps generated token to ersponse.
func (h *AuthHandler) respondWithAuthToken(c *gin.Context, token *core.Token) {
	c.JSON(http.StatusOK, response.AuthTokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn(),
		RefreshToken: token.RefreshToken,
	})
}
