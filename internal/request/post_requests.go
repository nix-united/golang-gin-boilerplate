package request

import validation "github.com/go-ozzo/ozzo-validation"

const (
	minPostTitleLength = 5
	maxPostTitleLength = 255
)

type BasicPost struct {
	Title   string `json:"title" binding:"required" example:"New Post"`
	Content string `json:"content" binding:"required" example:"Lorem Ipsum"`
}

func (p BasicPost) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Title, validation.Length(minPostTitleLength, maxPostTitleLength)),
		validation.Field(&p.Content, validation.Required),
	)
}

type CreatePostRequest struct {
	BasicPost
}

type UpdatePostRequest struct {
	BasicPost
}
