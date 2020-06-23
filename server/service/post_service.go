package service

import (
	"basic_server/server/model"
	"basic_server/server/repository"

	"github.com/jinzhu/gorm"
)

type PostService struct {
	PostRepository repository.PostRepository
	DB             *gorm.DB
}

func (service PostService) CreatePost(title, content string, userID uint) model.Post {
	return model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
}
