package utils

import "golang.org/x/crypto/bcrypt"

// Encryptor provides functions for strings encryption
type Encryptor interface {
	// Encrypt encrypts strings
	Encrypt(str string) (string, error)
}

type bcryptEncoder struct {
	cost int
}

// NewBcryptEncoder takes a cost value that is used
// during password encryption and returns a new Encryptor.
// The cost value can't be greater than 31.
// Otherwise, an error will be returned by Encrypt function
func NewBcryptEncoder(cost int) Encryptor {
	return bcryptEncoder{
		cost: cost,
	}
}

func (en bcryptEncoder) Encrypt(pass string) (string, error) {
	enPass, err := bcrypt.GenerateFromPassword(
		[]byte(pass),
		en.cost,
	)

	if err != nil {
		return "", err
	}

	return string(enPass), nil
}
