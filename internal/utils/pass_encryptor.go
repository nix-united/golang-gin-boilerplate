package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptEncoder struct {
	cost int
}

func NewBcryptEncoder(cost int) BcryptEncoder {
	return BcryptEncoder{
		cost: cost,
	}
}

func (en BcryptEncoder) Encrypt(pass string) (string, error) {
	enPass, err := bcrypt.GenerateFromPassword(
		[]byte(pass),
		en.cost,
	)
	if err != nil {
		return "", fmt.Errorf("generate from password: %w", err)
	}

	return string(enPass), nil
}
