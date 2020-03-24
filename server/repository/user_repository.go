package repository

import (
	"basic_server/server/model"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (repository UserRepository) FindUserByEmail (email string) model.User {
	var user model.User
	repository.DB.Where("email = ?", email).Find(&user)

	return user
}