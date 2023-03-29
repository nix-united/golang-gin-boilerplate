package request

import validation "github.com/go-ozzo/ozzo-validation"

type RegisterRequest struct {
	*BasicAuthRequest
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
}

func (rr *RegisterRequest) Validate() error {
	err := rr.BasicAuthRequest.Validate()
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&rr,
		validation.Field(&rr.FullName, validation.Required),
	)
}
