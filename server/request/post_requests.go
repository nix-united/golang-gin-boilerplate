package request

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"  example:"New Post"`
	Content string `json:"content" binding:"required"  example:"Lorem Ipsum"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" bind:"required"  example:"Title 1"`
	Content string `json:"content" bind:"required"  example:"description"`
}
