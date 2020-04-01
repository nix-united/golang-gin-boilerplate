package repository

import (
	"basic_server/server/model"
	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func (repository PostRepository) GetAll() []model.Post {
	var posts []model.Post
	repository.DB.Find(&posts)

	return posts
}

func (repository PostRepository) GetById(id uint) model.Post {
	var post model.Post
	repository.DB.First(&post, id)

	return post
}

func (repository PostRepository) Save(post model.Post) model.Post {
	repository.DB.Save(&post)

	return post
}

func (repository PostRepository) Delete(post model.Post) {
	repository.DB.Delete(&post)
}
