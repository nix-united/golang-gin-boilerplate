package provider

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/repository"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const identityKey = "id"

type Success struct {
	Code   int    `json:"code" example:"200"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

var once sync.Once

var mw *jwtAuthMiddleware

func NewJwtAuth(db *gorm.DB) JwtAuthMiddleware {
	once.Do(func() {
		var err error

		mw = &jwtAuthMiddleware{
			databaseDriver: db,
		}

		mw.authMiddleware, err = jwt.New(mw.prepareMiddleware())

		if err != nil {
			log.Fatal("JWT error")
		}
	})

	return mw
}

type JwtAuthMiddleware interface {
	Middleware() *jwt.GinJWTMiddleware
	Refresh(c *gin.Context)
}

type jwtAuthMiddleware struct {
	databaseDriver *gorm.DB
	authMiddleware *jwt.GinJWTMiddleware
}

func (mw *jwtAuthMiddleware) Middleware() *jwt.GinJWTMiddleware {
	return mw.authMiddleware
}

func (mw *jwtAuthMiddleware) prepareMiddleware() *jwt.GinJWTMiddleware {
	jwtSettings, err := utils.NewJwtEnvVars()

	if err != nil {
		log.Fatal(err)
	}

	middleware := &jwt.GinJWTMiddleware{
		Realm:                 jwtSettings.Realm(),
		Key:                   []byte(jwtSettings.Secret()),
		Timeout:               jwtSettings.Expiration(),
		MaxRefresh:            jwtSettings.RefreshTime(),
		IdentityKey:           identityKey,
		PayloadFunc:           addUserIDToClaims,
		IdentityHandler:       extractIdentityKeyFromClaims,
		Authorizator:          mw.isUserValid,
		Authenticator:         mw.authenticate,
		HTTPStatusMessageFunc: takeAppropriateErrorMessage,
		TimeFunc:              time.Now,
	}

	return middleware
}

// authenticate godoc
// @Summary Authenticate a user
// @Description Perform user login
// @ID user-login
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body request.BasicAuthRequest true "User's credentials"
// @Success 200 {object} Success
// @Failure 401 {object} response.Error
// @Router /login [post]
func (mw jwtAuthMiddleware) authenticate(c *gin.Context) (interface{}, error) {
	var authRequest request.BasicAuthRequest

	if err := c.ShouldBindJSON(&authRequest); err != nil {
		return model.User{}, jwt.ErrMissingLoginValues
	}

	userRepository := repository.NewUserRepository(mw.databaseDriver)

	user, err := userRepository.GetByEmail(c.Request.Context(), authRequest.Email)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password)) != nil) {
		return user, jwt.ErrFailedAuthentication
	}

	return user, nil
}

// refresh godoc
// @Summary Refresh token
// @Description Refresh user's token
// @ID refresh-token
// @Tags User Actions
// @Produce json
// @Success 200 {object} Success
// @Failure 401 {object} response.Error
// @Security ApiKeyAuth
// @Router /refresh [get]
func (mw jwtAuthMiddleware) Refresh(c *gin.Context) {
	mw.Middleware().RefreshHandler(c)
}

func (mw jwtAuthMiddleware) isUserValid(data interface{}, c *gin.Context) bool {
	userID, ok := data.(float64)
	if !ok {
		return false
	}

	userRepository := repository.NewUserRepository(mw.databaseDriver)

	_, err := userRepository.GetByID(c.Request.Context(), uint(userID))
	if err != nil {
		slog.WarnContext(c.Request.Context(), "Failed to get user by ID to authorize", "err", err)

		return false
	}

	return true
}

func extractIdentityKeyFromClaims(c *gin.Context) interface{} {
	identity, ok := jwt.ExtractClaims(c)[identityKey].(float64)
	if !ok {
		return 0
	}

	return identity
}

func addUserIDToClaims(data interface{}) jwt.MapClaims {
	if user, ok := data.(*model.User); ok {
		return jwt.MapClaims{
			identityKey: user.ID,
		}
	}

	return jwt.MapClaims{}
}

func takeAppropriateErrorMessage(err error, _ *gin.Context) string {
	switch err {
	case jwt.ErrMissingLoginValues:
		return "Email and password are required"
	case jwt.ErrFailedAuthentication:
		return "Invalid email or password"
	}

	return err.Error()
}
