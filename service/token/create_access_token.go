package token

import (
	"basic_server/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (tokenService *JWTTokenService) CreateAccessToken(user *model.User) (accessToken string, exp int64, err error) {
	exp = time.Now().Add(time.Hour * ExpireCount).Unix()
	claims := &JwtCustomClaims{
		user.FullName,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(tokenService.Config.Auth.AccessSecret))
	if err != nil {
		return "", 0, err
	}

	return t, exp, err
}
