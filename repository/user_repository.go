package repository

import (
	"basic_server/model"
	"errors"

	"github.com/jinzhu/gorm"
)

type UserRepositoryI interface {
	FindUserByEmail(email string) (model.User, error)
	FindUserByID(ID int) model.User
	StoreUser(user model.User) error
}

type UserRepository struct {
	storage *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryI {
	return &UserRepository{storage: db}
}

func (repo *UserRepository) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	err := repo.storage.Where("email = ?", email).Find(&user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}
	return user, nil
}

func (repo *UserRepository) FindUserByID(id int) model.User {
	var user model.User
	repo.storage.Where("id = ?", id).First(&user)

	return user
}

func (repo *UserRepository) StoreUser(user model.User) error { //nolint
	return repo.storage.Create(&user).Error
}
