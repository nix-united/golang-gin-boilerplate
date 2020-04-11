package provider

import (
	"basic_server/server/model"
	"basic_server/server/repository"
	"basic_server/server/request"
	"basic_server/server/service"
	"log"
	"sync"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const identityKey = "user_id"

var once sync.Once

var mw *jwtAuthMiddleware

func NewJwtAuth(db *gorm.DB) *jwtAuthMiddleware {
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

type jwtAuthMiddleware struct {
	databaseDriver *gorm.DB
	authMiddleware *jwt.GinJWTMiddleware
}

func (mw *jwtAuthMiddleware) Middleware() *jwt.GinJWTMiddleware {
	return mw.authMiddleware
}

func (mw *jwtAuthMiddleware) prepareMiddleware() *jwt.GinJWTMiddleware {
	jwtSettings, err := service.NewJwtEnvVars()

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

func (mw jwtAuthMiddleware) authenticate(c *gin.Context) (interface{}, error) {
	var authRequest request.AuthRequest
	var user model.User

	if err := c.ShouldBind(&authRequest); err != nil {
		return user, jwt.ErrMissingLoginValues
	}

	userRepository := repository.UserRepository{DB: mw.databaseDriver}

	user = userRepository.FindUserByEmail(authRequest.Email)

	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password)) != nil) {
		return user, jwt.ErrFailedAuthentication
	}

	return user, nil
}

func (mw jwtAuthMiddleware) isUserValid(data interface{}, c *gin.Context) bool {
	userID, ok := data.(float64)

	if !ok {
		return false
	}

	userRepository := repository.UserRepository{DB: mw.databaseDriver}

	if userRepository.FindUserById(int(userID)).ID == 0 {
		return false
	}

	return true
}

func extractIdentityKeyFromClaims(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	return claims[identityKey].(float64)
}

func addUserIDToClaims(data interface{}) jwt.MapClaims {
	if user, ok := data.(model.User); ok {
		return jwt.MapClaims{
			identityKey: user.ID,
		}
	}

	return jwt.MapClaims{}
}

func takeAppropriateErrorMessage(err error, c *gin.Context) string {
	switch err {
	case jwt.ErrMissingLoginValues:
		return "Email and password are required"
	case jwt.ErrFailedAuthentication:
		return "Invalid email or password"
	}

	return err.Error()
}
