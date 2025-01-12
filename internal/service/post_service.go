package service

import (
	"net/http"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

type postRepository interface {
	GetAll(posts *[]model.Post) error
	GetByID(id int, post *model.Post) error
	Create(post *model.Post) error
	Save(post *model.Post) error
	Delete(post *model.Post) error
}

type PostService struct {
	postRepository postRepository
}

func NewPostService(postRepository postRepository) PostService {
	return PostService{postRepository: postRepository}
}

func (service PostService) GetAll(posts *[]model.Post) *RestError {
	if err := service.postRepository.GetAll(posts); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service PostService) GetByID(id int, post *model.Post) *RestError {
	if err := service.postRepository.GetByID(id, post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service PostService) Create(post *model.Post) *RestError {
	if err := service.postRepository.Create(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service PostService) Save(post *model.Post) *RestError {
	if err := service.postRepository.Save(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service PostService) Delete(post *model.Post) *RestError {
	if err := service.postRepository.Delete(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service PostService) CreatePost(title, content string, userID uint) (*model.Post, *RestError) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := service.postRepository.Create(post); err != nil {
		return nil, &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return post, nil
}
