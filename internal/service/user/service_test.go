package user_test

import (
	"errors"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type serviceMovkcs struct {
	userRepository  *MockuserRepository
	passwordService *MockpasswordService
}

func newService(t *testing.T) (*user.Service, serviceMovkcs) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	passwordService := NewMockpasswordService(ctrl)
	userService := user.NewService(userRepository, passwordService)

	mocks := serviceMovkcs{
		userRepository:  userRepository,
		passwordService: passwordService,
	}

	return userService, mocks
}

func TestService_CreateUser(t *testing.T) {
	registerRequest := request.RegisterRequest{
		BasicAuthRequest: request.BasicAuthRequest{
			Email:    "test@test.com",
			Password: "password",
		},
		FullName: "test full name",
	}

	storedUser := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "stored_user@test.com",
		Password: "stored encrypted password",
		FullName: "stored user full name",
	}

	expectedUserToCreate := &model.User{
		Email:    registerRequest.Email,
		Password: "encrypted password",
		FullName: registerRequest.FullName,
	}

	t.Run("It should propagate an error if failed to find user in database", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, errors.New("unkown db error"))

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorContains(t, err, "get user by email")
	})

	t.Run("It should propagate an error if failed to encrypt password", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, domain.ErrNotFound)

		mocks.passwordService.
			EXPECT().
			EncryptPassword("password").
			Return("", errors.New("encryption error"))

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorContains(t, err, "encrypt password")
	})

	t.Run("It should propagate an error if failed to store an user", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, domain.ErrNotFound)

		mocks.passwordService.
			EXPECT().
			EncryptPassword("password").
			Return("encrypted password", nil)

		mocks.userRepository.
			EXPECT().
			Create(gomock.Any(), expectedUserToCreate).
			Return(errors.New("store user error"))

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorContains(t, err, "store user")
	})

	t.Run("It should return an error if user already exists in database", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(storedUser, nil)

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})

	t.Run("It should create a new user", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, domain.ErrNotFound)

		mocks.passwordService.
			EXPECT().
			EncryptPassword("password").
			Return("encrypted password", nil)

		mocks.userRepository.
			EXPECT().
			Create(gomock.Any(), expectedUserToCreate).
			Return(nil)

		err := service.CreateUser(t.Context(), registerRequest)
		assert.NoError(t, err)
	})
}

func TestService_GetUserByEmail(t *testing.T) {
	storedUser := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "stored_user@test.com",
		Password: "stored encrypted password",
		FullName: "stored user full name",
	}

	t.Run("It should fetch user by email from the repository", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), storedUser.Email).
			Return(storedUser, nil)

		actualUser, err := service.GetUserByEmail(t.Context(), storedUser.Email)
		require.NoError(t, err)
		assert.Equal(t, storedUser, actualUser)
	})

	t.Run("It should propagate error when failed to find such user in repository", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), storedUser.Email).
			Return(nil, domain.ErrNotFound)

		_, err := service.GetUserByEmail(t.Context(), storedUser.Email)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
