package repository

import (
	"errors"
	"fmt"

	operrors "github.com/nix-united/golang-gin-boilerplate/internal/errors"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return PostRepository{db: db}
}

func (r PostRepository) GetAll(posts *[]model.Post) error {
	return r.db.Find(posts).Error
}

func (r PostRepository) GetByID(id int) (*model.Post, error) {
	var post *model.Post
	err := r.db.Where("id = ? ", id).Take(post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Join(operrors.ErrPostNotFound, err)
	}
	if err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	return post, nil
}

func (r PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r PostRepository) Save(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r PostRepository) Delete(post *model.Post) error {
	return r.db.Delete(post).Error
}
