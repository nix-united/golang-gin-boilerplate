package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	cost int
}

func NewService(cost int) *Service {
	return &Service{cost: cost}
}

func (s *Service) VerifyPassword(actual, received string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(actual), []byte(received)); err != nil {
		return fmt.Errorf("compare hash and password: %w", err)
	}
	return nil
}

func (s *Service) EncryptPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("generate encrypted password: %w", err)
	}
	return string(encrypted), nil
}
