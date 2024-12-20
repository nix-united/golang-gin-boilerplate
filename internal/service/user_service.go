package service

import (
	"basic_server/internal/model"
	"basic_server/internal/request"
	"basic_server/internal/utils"
	"fmt"
)

//go:generate mockgen -source=$GOFILE -destination=user_service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRepository interface {
	FindUserByEmail(email string) (model.User, error)
	StoreUser(user model.User) error
}

// TODO: Change to private
type Encryptor interface {
	Encrypt(str string) (string, error)
}

// UserServiceI provides a use case level for the user entity
type UserServiceI interface {
	// CreateUser Create takes a request with new user credentials and registers it.
	// An error will be returned if a user exists in the system, or
	// if an error occurs during interaction with the database.
	CreateUser(req request.RegisterRequest, en utils.Encryptor) error
}

// UserService provides a use case level for the user entity
type UserService struct {
	userRepository userRepository
}

func NewUserService(userRepository userRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

// CreateUser Create takes a request with new user credentials and registers it.
// An error will be returned if a user exists in the system, or
// if an error occurs during interaction with the database.
func (s UserService) CreateUser(req request.RegisterRequest, en utils.Encryptor) error {
	user, err := s.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("find user by email: %w", err)
	}

	if user.ID != 0 {
		return NewErrUserAlreadyExists(
			"user already exist",
			"store a user",
		)
	}

	encryptedPassword, err := en.Encrypt(req.Password)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	err = s.userRepository.StoreUser(model.User{
		Email:    req.Email,
		Password: encryptedPassword,
		FullName: req.FullName,
	})
	if err != nil {
		return fmt.Errorf("store user: %w", err)
	}

	return nil
}
