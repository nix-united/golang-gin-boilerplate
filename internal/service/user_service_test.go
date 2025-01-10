package service_test

import (
	"errors"
	"testing"

	srverrors "github.com/nix-united/golang-gin-boilerplate/internal/errors"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type userServiceMocks struct {
	userRepository *MockuserRepository
	encryptor      *Mockencryptor
}

func newUserService(t *testing.T) (service.UserService, userServiceMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	encryptor := NewMockencryptor(ctrl)
	userService := service.NewUserService(userRepository, encryptor)

	mocks := userServiceMocks{
		userRepository: userRepository,
		encryptor:      encryptor,
	}

	return userService, mocks
}

func TestUserService_CreateUser(t *testing.T) {
	registerRequest := request.RegisterRequest{
		BasicAuthRequest: &request.BasicAuthRequest{
			Email:    "test@test.com",
			Password: "password",
		},
		FullName: "test full name",
	}

	storedUser := model.User{
		Email:    "test@test.com",
		Password: "encrypted password",
		FullName: "test full name",
	}

	userInDB := model.User{
		Model: gorm.Model{
			ID: 1,
		},
	}

	t.Run("It should propagate an error if failed to find user in database", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			FindUserByEmail("test@test.com").
			Return(model.User{}, errors.New("unkown db error"))

		err := service.CreateUser(registerRequest)
		assert.ErrorContains(t, err, "find user by email")
	})

	t.Run("It should propagate an error if failed to encrypt password", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			FindUserByEmail("test@test.com").
			Return(model.User{}, nil)

		mocks.encryptor.
			EXPECT().
			Encrypt("password").
			Return("", errors.New("encryption error"))

		err := service.CreateUser(registerRequest)
		assert.ErrorContains(t, err, "encrypt password")
	})

	t.Run("It should propagate an error if failed to store an user", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			FindUserByEmail("test@test.com").
			Return(model.User{}, nil)

		mocks.encryptor.
			EXPECT().
			Encrypt("password").
			Return("encrypted password", nil)

		mocks.userRepository.
			EXPECT().
			StoreUser(storedUser).
			Return(errors.New("store user error"))

		err := service.CreateUser(registerRequest)
		assert.ErrorContains(t, err, "store user")
	})

	t.Run("It should return an error if user already exists in database", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			FindUserByEmail("test@test.com").
			Return(userInDB, nil)

		err := service.CreateUser(registerRequest)

		var errInvalidStorageOperation srverrors.ErrInvalidStorageOperation
		require.ErrorAs(t, err, &errInvalidStorageOperation)

		assert.Equal(t, "user already exist", errInvalidStorageOperation.Error())
		assert.Equal(t, "store a user", errInvalidStorageOperation.Operation())
	})

	t.Run("It should create a new user", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			FindUserByEmail("test@test.com").
			Return(model.User{}, nil)

		mocks.encryptor.
			EXPECT().
			Encrypt("password").
			Return("encrypted password", nil)

		mocks.userRepository.
			EXPECT().
			StoreUser(storedUser).
			Return(nil)

		err := service.CreateUser(registerRequest)
		assert.NoError(t, err)
	})
}
