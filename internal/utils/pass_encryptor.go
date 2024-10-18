package utils

import "golang.org/x/crypto/bcrypt"

type Encryptor interface {
	Encrypt(str string) (string, error)
}

type bcryptEncoder struct {
	cost int
}

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
