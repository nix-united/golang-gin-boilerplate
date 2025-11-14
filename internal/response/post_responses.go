package response

import "github.com/nix-united/golang-gin-boilerplate/internal/model"

type PostResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewPostResponse(post *model.Post) PostResponse {
	return PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	}
}

func NewPostCollectionResponse(posts []model.Post) CollectionResponse[PostResponse] {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = NewPostResponse(&post)
	}

	return NewCollectionResponse(responses)
}
