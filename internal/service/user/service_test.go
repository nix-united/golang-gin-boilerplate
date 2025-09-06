package user_test

import (
	"errors"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/user"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type userServiceMocks struct {
	userRepository *MockuserRepository
	encryptor      *Mockencryptor
}

func newUserService(t *testing.T) (*user.Service, userServiceMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	encryptor := NewMockencryptor(ctrl)
	userService := user.NewService(userRepository, encryptor)

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

	storedUser := &model.User{
		Email:    "test@test.com",
		Password: "encrypted password",
		FullName: "test full name",
	}

	userInDB := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
	}

	t.Run("It should propagate an error if failed to find user in database", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, errors.New("unkown db error"))

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorContains(t, err, "get user by email")
	})

	t.Run("It should propagate an error if failed to encrypt password", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, domain.ErrNotFound)

		mocks.encryptor.
			EXPECT().
			Encrypt("password").
			Return("", errors.New("encryption error"))

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorContains(t, err, "encrypt password")
	})

	t.Run("It should propagate an error if failed to store an user", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, domain.ErrNotFound)

		mocks.encryptor.
			EXPECT().
			Encrypt("password").
			Return("encrypted password", nil)

		mocks.userRepository.
			EXPECT().
			Create(gomock.Any(), storedUser).
			Return(errors.New("store user error"))

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorContains(t, err, "store user")
	})

	t.Run("It should return an error if user already exists in database", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(userInDB, nil)

		err := service.CreateUser(t.Context(), registerRequest)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})

	t.Run("It should create a new user", func(t *testing.T) {
		service, mocks := newUserService(t)

		mocks.userRepository.
			EXPECT().
			GetByEmail(gomock.Any(), "test@test.com").
			Return(nil, domain.ErrNotFound)

		mocks.encryptor.
			EXPECT().
			Encrypt("password").
			Return("encrypted password", nil)

		mocks.userRepository.
			EXPECT().
			Create(gomock.Any(), storedUser).
			Return(nil)

		err := service.CreateUser(t.Context(), registerRequest)
		assert.NoError(t, err)
	})
}
