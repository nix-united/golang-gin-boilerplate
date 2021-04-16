package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	minPathLength = 8
)

type RefreshRequest struct {
	Token string `json:"token" validate:"required" example:"refresh_token"`
}

type BasicAuthRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"11111111"`
}

func (ar *BasicAuthRequest) Validate() error {
	return validation.ValidateStruct(ar,
		validation.Field(ar.Email, is.Email),
		validation.Field(ar.Password, validation.Length(minPathLength, 0)),
	)
}
