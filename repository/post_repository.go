package repository

import (
	"basic_server/model"

	"gorm.io/gorm"
)

type PostRepositoryI interface {
	GetAll(posts *[]model.Post) error
	GetByID(id int, post *model.Post) error
	Create(post *model.Post) error
	Save(post *model.Post) error
	Delete(post *model.Post) error
}

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}

func (repository *PostRepository) GetAll(posts *[]model.Post) error {
	return repository.DB.Find(posts).Error
}

func (repository *PostRepository) GetByID(id int, post *model.Post) error {
	return repository.DB.Where("id = ? ", id).Find(post).Error
}

func (repository *PostRepository) Create(post *model.Post) error {
	return repository.DB.Create(post).Error
}

func (repository *PostRepository) Save(post *model.Post) error {
	return repository.DB.Save(post).Error
}

func (repository *PostRepository) Delete(post *model.Post) error {
	return repository.DB.Delete(post).Error
}
