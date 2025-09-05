package integration_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	userRespository := repository.NewUserRepository(gormDB)

	user := &model.User{
		Email:    "test_user_repository_user@example.com",
		Password: "test_user_repository_user_password",
		FullName: "test_user_repository_user_full_name",
	}

	t.Run("It should create a user", func(t *testing.T) {
		err := userRespository.Create(t.Context(), user)
		require.NoError(t, err)
	})

	t.Run("It should get user by ID", func(t *testing.T) {
		gotUser, err := userRespository.GetByID(t.Context(), user.ID)
		require.NoError(t, err)

		assert.Equal(t, user.Email, gotUser.Email)
	})

	t.Run("It should return ErrNotFound error when user with such ID not found", func(t *testing.T) {
		_, err := userRespository.GetByID(t.Context(), 999999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("It should get user by email", func(t *testing.T) {
		gotUser, err := userRespository.GetByEmail(t.Context(), user.Email)
		require.NoError(t, err)

		assert.Equal(t, user.Email, gotUser.Email)
	})

	t.Run("It should return ErrNotFound error when user with such email not found", func(t *testing.T) {
		_, err := userRespository.GetByEmail(t.Context(), "unknown@email.com")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
