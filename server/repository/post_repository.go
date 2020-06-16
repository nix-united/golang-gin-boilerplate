package repository

import (
	"basic_server/server/model"

	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func (repository PostRepository) GetAll(posts *[]model.Post) {
	repository.DB.Find(posts)
}

func (repository PostRepository) GetByID(id int, post *model.Post) {
	repository.DB.Where("id = ? ", id).Find(post)
}

func (repository PostRepository) Create(post *model.Post) {
	repository.DB.Create(post)
}

func (repository PostRepository) Save(post *model.Post) {
	repository.DB.Save(post)
}

func (repository PostRepository) Delete(post *model.Post) {
	repository.DB.Delete(post)
}
