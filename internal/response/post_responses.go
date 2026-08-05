package response

import (
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

type PostResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewPostResponse(post *model.Post) PostResponse {
	return PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}
}

type PostCollectionResponse = CollectionResponse[PostResponse]

func NewPostCollectionResponse(posts []model.Post, total int64, offset, limit int) PostCollectionResponse {
	responses := make([]PostResponse, len(posts))
	for i := range posts {
		responses[i] = NewPostResponse(&posts[i])
	}

	return NewCollectionResponse(responses, total, offset, limit)
}
