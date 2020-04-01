package service

import (
	"basic_server/server/model"
	"basic_server/server/repository"
)

type PostService struct {
	PostRepository repository.PostRepository
}


func (service PostService) CreatePost(title, content string) model.Post {
	return model.Post{
		Title:   title,
		Content: content,
	}
}

func (service PostService) GetAll() []model.Post {
	return service.PostRepository.GetAll()
}

func (service PostService) GetById(id uint) model.Post {
	return service.PostRepository.GetById(id)
}

func (service PostService) Save(post model.Post) model.Post {
	service.PostRepository.Save(post)

	return post
}

func (service PostService) Delete(post model.Post) {
	service.PostRepository.Delete(post)
}
