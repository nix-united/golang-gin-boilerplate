package service

import (
	"net/http"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/repository"
)

type PostServiceI interface {
	CreatePost(title, content string, userID uint) (*model.Post, *RestError)
	GetAll(posts *[]model.Post) *RestError
	GetByID(id int, post *model.Post) *RestError
	Create(post *model.Post) *RestError
	Save(post *model.Post) *RestError
	Delete(post *model.Post) *RestError
}

type PostService struct {
	PostRepository repository.PostRepositoryI
}

func (service *PostService) GetAll(posts *[]model.Post) *RestError {
	if err := service.PostRepository.GetAll(posts); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service *PostService) GetByID(id int, post *model.Post) *RestError {
	if err := service.PostRepository.GetByID(id, post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service *PostService) Create(post *model.Post) *RestError {
	if err := service.PostRepository.Create(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service *PostService) Save(post *model.Post) *RestError {
	if err := service.PostRepository.Save(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service *PostService) Delete(post *model.Post) *RestError {
	if err := service.PostRepository.Delete(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return nil
}

func (service *PostService) CreatePost(title, content string, userID uint) (*model.Post, *RestError) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := service.PostRepository.Create(post); err != nil {
		return nil, &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return post, nil
}

func NewPostService(postRepo repository.PostRepositoryI) PostServiceI {
	return &PostService{PostRepository: postRepo}
}
