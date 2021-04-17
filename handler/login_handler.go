package handler

import (
	"basic_server/repository"
	"basic_server/request"
	"basic_server/response"
	"basic_server/service/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginHandler struct {
	UserRepository repository.UserRepositoryI
	TokenService   token.JWTTokenService
}

func NewLoginHandler(userRepository repository.UserRepositoryI) *LoginHandler {
	return &LoginHandler{
		UserRepository: userRepository,
	}
}

// Login godoc
// @Summary Sign In user
// @Description Sign In User to the system
// @ID login-user
// @Tags Auth Action
// @Accept json
// @Produce json
// @Param params body request.BasicAuthRequest true "Post user login form"
// @Success 200 {object} response.LoginResponse
// @Failure 401 {object} response.Error
// @Security ApiKeyAuth
// @Router /auth/login [post]
func (h *LoginHandler) Login(context *gin.Context) {
	var authRequest request.BasicAuthRequest

	if err := context.ShouldBind(&authRequest); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		return
	}

	if err := authRequest.Validate(); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty or not valid")
		return
	}

	user, err := h.UserRepository.FindUserByEmail(authRequest.Email)
	if err != nil ||
		user.ID == 0 ||
		(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password)) != nil) {
		response.ErrorResponse(context, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, exp, err := h.TokenService.CreateAccessToken(&user)
	if err != nil {
		response.ErrorResponse(context, http.StatusUnauthorized, "Invalid or expired session")
		return
	}
	refreshToken, err := h.TokenService.CreateRefreshToken(&user)
	if err != nil {
		response.ErrorResponse(context, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	response.Response(
		context,
		http.StatusOK,
		&response.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Exp:          exp,
		},
	)
}
