package request

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterRequest struct {
	BasicAuthRequest

	FullName string `json:"full_name" validate:"required" example:"John Doe"`
}

func (r RegisterRequest) Validate() error {
	return errors.Join(
		r.BasicAuthRequest.Validate(),
		validation.ValidateStruct(&r, validation.Field(&r.FullName, validation.Required)),
	)
}
