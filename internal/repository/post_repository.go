package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return nil, fmt.Errorf("execute insert post query: %w", err)
	}

	return post, nil
}

func (r *PostRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Post{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("execute count number of posts query: %w", err)
	}

	return count, nil
}

func (r *PostRepository) List(ctx context.Context, filters domain.PostFilters) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.WithContext(ctx).
		Offset(int(filters.Offset)).
		Limit(int(filters.Limit)).
		Find(&posts).
		Error
	if err != nil {
		return nil, fmt.Errorf("execute select posts query: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	var post *model.Post
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(domain.ErrNotFound, err)
		}

		return nil, fmt.Errorf("execute select post by id query: %w", err)
	}

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, post *model.Post) error {
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
