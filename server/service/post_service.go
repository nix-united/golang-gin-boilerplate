package service

import (
	"basic_server/server/model"
)

type PostService struct{}

func (service PostService) CreatePost(title, content string) model.Post {
	return model.Post{
		Title:   title,
		Content: content,
	}
}
