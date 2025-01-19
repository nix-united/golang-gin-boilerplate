package service

import (
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=post_service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postRepository interface {
	GetAll(posts *[]model.Post) error
	GetByID(id int) (*model.Post, error)
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

func (s PostService) CreatePost(title, content string, userID uint) (*model.Post, error) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := s.postRepository.Create(post); err != nil {
		return nil, fmt.Errorf("create post in repository: %w", err)
	}

	return post, nil
}

func (s PostService) Create(post *model.Post) error {
	if err := s.postRepository.Create(post); err != nil {
		return fmt.Errorf("create post in repository: %w", err)
	}

	return nil
}

func (s PostService) GetAll(posts *[]model.Post) error {
	if err := s.postRepository.GetAll(posts); err != nil {
		return fmt.Errorf("get all posts from repository: %w", err)
	}

	return nil
}

func (s PostService) GetByID(id int) (*model.Post, error) {
	post, err := s.postRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get post by id from repository: %w", err)
	}

	return post, nil
}

func (s PostService) Save(post *model.Post) error {
	if err := s.postRepository.Save(post); err != nil {
		return fmt.Errorf("save post in repository: %w", err)
	}

	return nil
}

func (s PostService) Delete(post *model.Post) error {
	if err := s.postRepository.Delete(post); err != nil {
		return fmt.Errorf("delete post from repository: %w", err)
	}

	return nil
}
