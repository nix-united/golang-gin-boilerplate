package handler

import (
	"errors"
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/ccoveille/go-safecast"
	"github.com/gin-gonic/gin"
)

// getUserIDFromContext extracts the user ID from JWT claims set by
// the authorization middleware.
func getUserIDFromContext(c *gin.Context) (uint, error) {
	claims := jwt.ExtractClaims(c)

	id, ok := claims[identityKey].(float64)
	if !ok {
		return 0, errors.New("missing user id in claims")
	}

	parsed, err := safecast.Convert[uint](id)
	if err != nil {
		return 0, fmt.Errorf("convert user id to uint: %w", err)
	}

	return parsed, nil
}
