package repository

import (
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

func (r PostRepository) GetByID(id int, post *model.Post) error {
	return r.db.Where("id = ? ", id).Find(post).Error
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
