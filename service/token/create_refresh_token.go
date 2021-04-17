package token

import (
	"basic_server/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (tokenService *JWTTokenService) CreateRefreshToken(user *model.User) (t string, err error) {
	claimsRefresh := &JwtCustomRefreshClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * ExpireRefreshCount).Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	rt, err := refreshToken.SignedString([]byte(tokenService.Config.Auth.RefreshSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}
