package handler

import (
	"errors"
	"net/http"

	operror "github.com/nix-united/golang-gin-boilerplate/internal/errors"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=$GOFILE -destination=auth_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userService interface {
	CreateUser(req request.RegisterRequest) error
}

type AuthHandler struct {
	userService userService
}

func NewAuthHandler(userService userService) AuthHandler {
	return AuthHandler{
		userService: userService,
	}
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
func (h AuthHandler) RegisterUser(c *gin.Context) {
	var registerRequest request.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		response.ErrorResponse(
			c,
			http.StatusUnprocessableEntity,
			"Required fields are empty or email is not valid",
		)
		return
	}

	if err := registerRequest.Validate(); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	err := h.userService.CreateUser(registerRequest)
	if err != nil {
		var errOperation operror.ErrInvalidStorageOperation
		if errors.As(err, &errOperation) {
			response.ErrorResponse(
				c,
				http.StatusUnprocessableEntity,
				errOperation.Error(),
			)
			return
		}

		response.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Oops, something went wrong...",
		)
		return
	}

	response.SuccessResponse(c, "Successfully registered")
}
