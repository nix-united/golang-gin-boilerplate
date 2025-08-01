package repository

import (
	"context"
	"errors"
	"fmt"

	operrors "github.com/nix-united/golang-gin-boilerplate/internal/errors"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return fmt.Errorf("execute insert post query: %w", err)
	}

	return nil
}

func (r *PostRepository) GetAll(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("execute select posts query: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) GetByID(ctx context.Context, id int) (*model.Post, error) {
	var post *model.Post
	err := r.db.WithContext(ctx).Where("id = ? ", id).Take(post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Join(operrors.ErrPostNotFound, err)
	}
	if err != nil {
		return nil, fmt.Errorf("execute select post by id query: %w", err)
	}

	return post, nil
}

func (r *PostRepository) Save(ctx context.Context, post *model.Post) error {
	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		return fmt.Errorf("execute update post query: %w", err)
	}

	return nil
}

func (r *PostRepository) Delete(ctx context.Context, post *model.Post) error {
	if err := r.db.WithContext(ctx).Delete(post).Error; err != nil {
		return fmt.Errorf("execute delete post query: %w", err)
	}

	return nil
}
