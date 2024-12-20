package request

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterRequest struct {
	*BasicAuthRequest
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
}

func (rr *RegisterRequest) Validate() error {
	return errors.Join(
		rr.BasicAuthRequest.Validate(),
		validation.ValidateStruct(rr, validation.Field(&rr.FullName, validation.Required)),
	)
}
