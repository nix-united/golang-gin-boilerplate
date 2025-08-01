package user

import (
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRepository interface {
	FindUserByEmail(email string) (model.User, error)
	StoreUser(user model.User) error
}

type encryptor interface {
	Encrypt(str string) (string, error)
}

// Service provides a use case level for the user entity
type Service struct {
	userRepository userRepository
	encryptor      encryptor
}

func NewService(userRepository userRepository, enencryptor encryptor) Service {
	return Service{
		userRepository: userRepository,
		encryptor:      enencryptor,
	}
}

// CreateUser Create takes a request with new user credentials and registers it.
// An error will be returned if a user exists in the system, or
// if an error occurs during interaction with the database.
func (s Service) CreateUser(req request.RegisterRequest) error {
	user, err := s.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("find user by email: %w", err)
	}

	if user.ID != 0 {
		return service.NewErrUserAlreadyExists(
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
