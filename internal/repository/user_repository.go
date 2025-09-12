package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("execute insert user query: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(domain.ErrNotFound, err)
		}

		return nil, fmt.Errorf("execute select user by id query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(domain.ErrNotFound, err)
		}

		return nil, fmt.Errorf("execute select user by email query: %w", err)
	}

	return user, nil
}
