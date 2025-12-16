package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const minPasswordLength = 8

type RefreshRequest struct {
	Token string `json:"token" validate:"required" example:"refresh_token"`
}

type BasicAuthRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"11111111"`
}

func (r BasicAuthRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Email, is.Email),
		validation.Field(&r.Password, validation.Length(minPasswordLength, 0)),
	)
}
