package password_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/service/password"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestService(t *testing.T) {
	service := password.NewService(bcrypt.DefaultCost)

	const testPassword = "test-password"

	encryptedPassword, err := service.EncryptPassword(testPassword)
	require.NoError(t, err)

	err = service.VerifyPassword(encryptedPassword, testPassword)
	require.NoError(t, err)
}
