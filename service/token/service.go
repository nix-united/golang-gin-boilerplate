package token

import (
	"basic_server/config"
	"basic_server/model"
	"github.com/dgrijalva/jwt-go"
)

const ExpireCount = 2
const ExpireRefreshCount = 168

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

type ServiceWrapper interface {
	CreateAccessToken(user *model.User) (accessToken string, exp int64, err error)
	CreateRefreshToken(user *model.User) (t string, err error)
}

type JWTTokenService struct {
	Config *config.Config
}

func NewTokenService(cfg *config.Config) *JWTTokenService {
	return &JWTTokenService{
		Config: cfg,
	}
}
