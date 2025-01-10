package service

import (
	"basic_server/internal/model"
	"basic_server/internal/request"
	"fmt"
)

//go:generate mockgen -source=$GOFILE -destination=user_service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRepository interface {
	FindUserByEmail(email string) (model.User, error)
	StoreUser(user model.User) error
}

type encryptor interface {
	Encrypt(str string) (string, error)
}

// UserService provides a use case level for the user entity
type UserService struct {
	userRepository userRepository
	encryptor      encryptor
}

func NewUserService(userRepository userRepository, enencryptor encryptor) UserService {
	return UserService{
		userRepository: userRepository,
		encryptor:      enencryptor,
	}
}

// CreateUser Create takes a request with new user credentials and registers it.
// An error will be returned if a user exists in the system, or
// if an error occurs during interaction with the database.
func (s UserService) CreateUser(req request.RegisterRequest) error {
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

	encryptedPassword, err := s.encryptor.Encrypt(req.Password)
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
