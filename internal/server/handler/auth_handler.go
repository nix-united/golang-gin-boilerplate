package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=$GOFILE -destination=auth_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userService interface {
	CreateUser(ctx context.Context, registerRequest request.RegisterRequest) error
}

type AuthHandler struct {
	userService userService
}

func NewAuthHandler(userService userService) *AuthHandler {
	return &AuthHandler{userService: userService}
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
// @Failure 422 {object} response.Error
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
