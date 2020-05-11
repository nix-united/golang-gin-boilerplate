package service

import (
	"errors"
	"os"
	"strconv"
	"time"
)

var (
	errSecretKeyIsNotSet = errors.New("jwt secret key is not set")

	errRealmIsNotSet = errors.New("jwt realm is not set")

	errExpirationTimeHasNotBeenLoaded = errors.New("an error has occurred during jwt expiration time loading")

	errMaxRefreshTimeHasNotBeenLoadedE = errors.New("an error has occurred during jwt max refresh time loading")
)

func NewJwtEnvVars() (JwtEnvVars, error) {
	var jwtVars *jwtEnvVars
	var jwtSecret string
	var jwtRealm string
	var jwtExpration int
	var jwtMaxRefreshTime int
	var err error

	if jwtSecret = os.Getenv("JWT_SECRET"); jwtSecret == "" {
		return jwtVars, errSecretKeyIsNotSet
	}

	if jwtRealm = os.Getenv("JWT_REALM"); jwtRealm == "" {
		return jwtVars, errRealmIsNotSet
	}

	if jwtExpration, err = strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME")); err != nil {
		return jwtVars, errExpirationTimeHasNotBeenLoaded
	}

	if jwtMaxRefreshTime, err = strconv.Atoi(os.Getenv("JWT_REFRESH_TIME")); err != nil {
		return jwtVars, errMaxRefreshTimeHasNotBeenLoadedE
	}

	return &jwtEnvVars{
		secret:         jwtSecret,
		realm:          jwtRealm,
		expirationTime: time.Duration(jwtExpration) * time.Second,
		maxRefreshTime: time.Duration(jwtMaxRefreshTime) * time.Second,
	}, nil
}

type JwtEnvVars interface {
	Secret() string
	Realm() string
	Expiration() time.Duration
	RefreshTime() time.Duration
}

type jwtEnvVars struct {
	secret         string
	realm          string
	expirationTime time.Duration
	maxRefreshTime time.Duration
}

func (jwt *jwtEnvVars) Secret() string {
	return jwt.secret
}

func (jwt *jwtEnvVars) Realm() string {
	return jwt.secret
}

func (jwt *jwtEnvVars) Expiration() time.Duration {
	return jwt.expirationTime
}

func (jwt *jwtEnvVars) RefreshTime() time.Duration {
	return jwt.maxRefreshTime
}
