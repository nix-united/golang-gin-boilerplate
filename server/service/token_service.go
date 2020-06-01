package service

import (
	"basic_server/server/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const ExpireCount = 72

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	jwt.StandardClaims
}

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (tokenService *TokenService) CreateToken(user *model.User) (string, error) {
	claims := &JwtCustomClaims{
		user.FullName,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * ExpireCount).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, err
}


