package post

import (
	"context"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postRepository interface {
	GetAll(ctx context.Context) ([]model.Post, error)
	GetByID(ctx context.Context, id int) (*model.Post, error)
	Create(ctx context.Context, post *model.Post) error
	Save(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, post *model.Post) error
}

type Service struct {
	postRepository postRepository
}

func NewService(postRepository postRepository) *Service {
	return &Service{postRepository: postRepository}
}

func (s *Service) CreatePost(ctx context.Context, title, content string, userID uint) (*model.Post, error) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := s.postRepository.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("create post in repository: %w", err)
	}

	return post, nil
}

func (s *Service) Create(ctx context.Context, post *model.Post) error {
	if err := s.postRepository.Create(ctx, post); err != nil {
		return fmt.Errorf("create post in repository: %w", err)
	}

	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]model.Post, error) {
	posts, err := s.postRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all posts from repository: %w", err)
	}

	return posts, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*model.Post, error) {
	post, err := s.postRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post by id from repository: %w", err)
	}

	return post, nil
}

func (s *Service) Save(ctx context.Context, post *model.Post) error {
	if err := s.postRepository.Save(ctx, post); err != nil {
		return fmt.Errorf("save post in repository: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, post *model.Post) error {
	if err := s.postRepository.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete post from repository: %w", err)
	}

	return nil
}
