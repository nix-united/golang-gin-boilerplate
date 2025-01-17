package service

import (
	"net/http"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=post_service_mock_test.go -package=${GOPACKAGE}_test -typed=true

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

func (s PostService) CreatePost(title, content string, userID uint) (*model.Post, *RestError) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := s.postRepository.Create(post); err != nil {
		return nil, &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return post, nil
}

func (s PostService) Create(post *model.Post) *RestError {
	if err := s.postRepository.Create(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}

func (s PostService) GetAll(posts *[]model.Post) *RestError {
	if err := s.postRepository.GetAll(posts); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}

func (s PostService) GetByID(id int, post *model.Post) *RestError {
	if err := s.postRepository.GetByID(id, post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}

func (s PostService) Save(post *model.Post) *RestError {
	if err := s.postRepository.Save(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}

func (s PostService) Delete(post *model.Post) *RestError {
	if err := s.postRepository.Delete(post); err != nil {
		return &RestError{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}
