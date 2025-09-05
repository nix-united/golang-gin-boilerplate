package post

import (
	"cmp"
	"context"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id uint) (*model.Post, error)
	List(ctx context.Context) ([]model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, post *model.Post) error
}

type Service struct {
	postRepository postRepository
}

func NewService(postRepository postRepository) *Service {
	return &Service{postRepository: postRepository}
}

func (s *Service) Create(ctx context.Context, userID uint, title, content string) (*model.Post, error) {
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

func (s *Service) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	post, err := s.postRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post by id from repository: %w", err)
	}

	return post, nil
}

func (s *Service) List(ctx context.Context) ([]model.Post, error) {
	posts, err := s.postRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all posts from repository: %w", err)
	}

	return posts, nil
}

func (s *Service) UpdateByUser(ctx context.Context, userID, postID uint, title, content string) (*model.Post, error) {
	post, err := s.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	if post.UserID != userID {
		return nil, fmt.Errorf("post belongs to a different user: %w", domain.ErrForbidden)
	}

	post.Title = cmp.Or(title, post.Title)
	post.Content = cmp.Or(content, post.Content)

	if err := s.postRepository.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("Update post in repository: %w", err)
	}

	return post, nil
}

func (s *Service) DeleteByUser(ctx context.Context, userID, postID uint) error {
	post, err := s.postRepository.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post by id: %w", err)
	}

	if post.UserID != userID {
		return fmt.Errorf("post belongs to a different user: %w", domain.ErrForbidden)
	}

	if err := s.postRepository.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete post from repository: %w", err)
	}

	return nil
}
