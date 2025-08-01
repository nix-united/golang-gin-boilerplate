package repository

import (
	"context"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Find(user).Error
	if err != nil {
		return nil, fmt.Errorf("execute select user by email query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*model.User, error) {
	var user *model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("execute select user by id query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) StoreUser(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("execute insert user query: %w", err)
	}

	return nil
}
