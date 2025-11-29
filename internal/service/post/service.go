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
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, filters domain.PostFilters) ([]model.Post, error)
	GetByID(ctx context.Context, id uint) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, post *model.Post) error
}

type Service struct {
	postRepository postRepository
}

func NewService(postRepository postRepository) *Service {
	return &Service{postRepository: postRepository}
}

func (s *Service) Create(ctx context.Context, createPostRequest domain.CreatePostRequest) (*model.Post, error) {
	post := &model.Post{
		UserID:  createPostRequest.UserID,
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
	}

	post, err := s.postRepository.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("create post in repository: %w", err)
	}

	return post, nil
}

func (s *Service) Count(ctx context.Context) (int64, error) {
	count, err := s.postRepository.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("get posts count from repository: %w", err)
	}

	return count, nil
}

func (s *Service) List(ctx context.Context, filters domain.PostFilters) ([]model.Post, error) {
	posts, err := s.postRepository.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("get all posts from repository: %w", err)
	}

	return posts, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	post, err := s.postRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post by id from repository: %w", err)
	}

	return post, nil
}

func (s *Service) UpdateByUser(ctx context.Context, request domain.UpdatePostRequest) (*model.Post, error) {
	post, err := s.postRepository.GetByID(ctx, request.PostID)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	if post.UserID != request.UserID {
		return nil, fmt.Errorf("post belongs to a different user: %w", domain.ErrForbidden)
	}

	post.Title = cmp.Or(request.Title, post.Title)
	post.Content = cmp.Or(request.Content, post.Content)

	if err := s.postRepository.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("update post in repository: %w", err)
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
