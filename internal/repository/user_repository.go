package repository

import (
	"errors"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}

	return user, nil
}

func (r UserRepository) FindUserByID(id int) model.User {
	var user model.User
	r.db.Where("id = ?", id).First(&user)

	return user
}

func (r UserRepository) StoreUser(user model.User) error {
	return r.db.Create(&user).Error
}
