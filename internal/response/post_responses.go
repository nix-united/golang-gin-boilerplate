package response

import (
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

type PostResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewPostResponse(post *model.Post) PostResponse {
	return PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}
}

func NewPostCollectionResponse(posts []model.Post, total, offset, limit int64) CollectionResponse[PostResponse] {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = NewPostResponse(&post)
	}
	return NewCollectionResponse(responses, total, offset, limit)
}
