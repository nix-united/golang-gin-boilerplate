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