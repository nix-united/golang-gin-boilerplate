package handler

import (
	"basic_server/server/repository"
	"basic_server/server/request"
	"basic_server/server/response"
	"basic_server/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	DB *gorm.DB
}

// Register godoc
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
func (handler *RegisterHandler) Register() gin.HandlerFunc {
	return func(context *gin.Context) {
		var registerRequest request.RegisterRequest
		var newUserService service.NewUserService

		err := context.ShouldBind(&registerRequest)

		if err != nil {
			response.ErrorResponse(
				context,
				http.StatusUnprocessableEntity,
				"Required fields are empty or email is not valid",
			)
			return
		}

		userRepository := repository.UserRepository{DB: handler.DB}

		if userRepository.FindUserByEmail(registerRequest.Email).ID != 0 {
			response.ErrorResponse(
				context,
				http.StatusUnprocessableEntity,
				"User already exists",
			)
			return
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(registerRequest.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			context.Status(http.StatusInternalServerError)
			return
		}

		newUser := newUserService.CreateUser(
			registerRequest.Email,
			string(encryptedPassword),
			registerRequest.FullName,
		)

		handler.DB.Create(&newUser)

		response.SuccessResponse(context, "Successfully registered")
	}
}
