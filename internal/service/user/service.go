package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
)

//go:generate mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type passwordService interface {
	EncryptPassword(password string) (string, error)
}

// Service provides a use case level for the user entity
type Service struct {
	userRepository  userRepository
	passwordService passwordService
}

func NewService(userRepository userRepository, passwordService passwordService) *Service {
	return &Service{
		userRepository:  userRepository,
		passwordService: passwordService,
	}
}

// CreateUser Create takes a request with new user credentials and registers it.
// An error will be returned if a user exists in the system, or
// if an error occurs during interaction with the database.
func (s *Service) CreateUser(ctx context.Context, req request.RegisterRequest) error {
	_, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("get user by email: %w", err)
	} else if err == nil {
		return domain.ErrAlreadyExists
	}

	encryptedPassword, err := s.passwordService.EncryptPassword(req.Password)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	err = s.userRepository.Create(ctx, &model.User{
		Email:    req.Email,
		Password: encryptedPassword,
		FullName: req.FullName,
	})
	if err != nil {
		return fmt.Errorf("store user: %w", err)
	}

	return nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email from repository: %w", err)
	}
	return user, nil
}
