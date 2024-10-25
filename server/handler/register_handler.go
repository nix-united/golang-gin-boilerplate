package handler

import (
	"net/http"

	operror "basic_server/internal/errors"
	"basic_server/internal/request"
	"basic_server/response"
	"basic_server/service"
	"basic_server/utils" //nolint

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct{} //nolint

func NewRegisterHandler() RegisterHandler { //nolint
	return RegisterHandler{}
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
func (handler RegisterHandler) RegisterUser(srv service.UserServiceI) gin.HandlerFunc {
	return func(context *gin.Context) {
		var registerRequest request.RegisterRequest
		var err error

		err = context.ShouldBind(&registerRequest)

		if err != nil {
			response.ErrorResponse(
				context,
				http.StatusUnprocessableEntity,
				"Required fields are empty or email is not valid",
			)
			return
		}

		err = srv.CreateUser(registerRequest, utils.NewBcryptEncoder(bcrypt.DefaultCost))

		if err != nil {
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
}
