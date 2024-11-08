package handler

import (
	operror "basic_server/internal/errors"
	"basic_server/internal/request"
	"basic_server/internal/response"
	"basic_server/internal/service"
	"basic_server/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserService service.UserServiceI
}

func NewAuthHandler(userService service.UserServiceI) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
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
func (h *AuthHandler) RegisterUser(context *gin.Context) {
	var (
		registerRequest request.RegisterRequest
		err             error
	)

	if err = context.ShouldBindJSON(&registerRequest); err != nil {
		response.ErrorResponse(
			context,
			http.StatusUnprocessableEntity,
			"Required fields are empty or email is not valid",
		)
		return
	}

	if err = h.UserService.CreateUser(registerRequest, utils.NewBcryptEncoder(bcrypt.DefaultCost)); err != nil {
		if operationErr, ok := err.(operror.ErrInvalidStorageOperation); ok {
			response.ErrorResponse(
				context,
				http.StatusUnprocessableEntity,
				operationErr.Error(),
			)
			return
		}

		response.ErrorResponse(
			context,
			http.StatusInternalServerError,
			"Oops, something went wrong...",
		)
		return
	}

	response.SuccessResponse(context, "Successfully registered")
}
