package request

import validation "github.com/go-ozzo/ozzo-validation"

type BasicPost struct {
	Title   string `json:"title" binding:"required"  example:"New Post"`
	Content string `json:"content" binding:"required"  example:"Lorem Ipsum"`
}

type CreatePostRequest struct {
	*BasicPost
}

type UpdatePostRequest struct {
	*BasicPost
}

func (bp *BasicPost) Validate() error {
	return validation.ValidateStruct(bp,
		validation.Field(bp.Title, validation.Required),
		validation.Field(bp.Content, validation.Required),
	)
}
