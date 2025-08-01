package post

import (
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postRepository interface {
	GetAll() ([]model.Post, error)
	GetByID(id int) (*model.Post, error)
	Create(post *model.Post) error
	Save(post *model.Post) error
	Delete(post *model.Post) error
}

type Service struct {
	postRepository postRepository
}

func NewService(postRepository postRepository) *Service {
	return &Service{postRepository: postRepository}
}

func (s *Service) CreatePost(title, content string, userID uint) (*model.Post, error) {
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

func (s *Service) Create(post *model.Post) error {
	if err := s.postRepository.Create(post); err != nil {
		return fmt.Errorf("create post in repository: %w", err)
	}

	return nil
}

func (s *Service) GetAll() ([]model.Post, error) {
	posts, err := s.postRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all posts from repository: %w", err)
	}

	return posts, nil
}

func (s *Service) GetByID(id int) (*model.Post, error) {
	post, err := s.postRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get post by id from repository: %w", err)
	}

	return post, nil
}

func (s *Service) Save(post *model.Post) error {
	if err := s.postRepository.Save(post); err != nil {
		return fmt.Errorf("save post in repository: %w", err)
	}

	return nil
}

func (s *Service) Delete(post *model.Post) error {
	if err := s.postRepository.Delete(post); err != nil {
		return fmt.Errorf("delete post from repository: %w", err)
	}

	return nil
}
