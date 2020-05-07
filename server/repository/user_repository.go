package repository

import (
	"basic_server/server/model"

	"github.com/jinzhu/gorm"
)

// UserRepository provides functionality for interacting with users storage
type UsersRepository interface {
	// FindByEmail takes an email and returns a user.
	// If a user isn't found, the function returns an empty user model
	FindUserByEmail(email string) (model.User, error)

	// FindById takes user id an returns a user.
	// If a user isn't found, the function returns an empty user model
	FindUserById(ID int) model.User

	// Store takes a user and saves it. The function returns
	// an error if it occurs during interaction with a storage
	StoreUser(user model.User) error
}

type usersRepository struct {
	storage *gorm.DB
}

// NewUserRepository returns an instance of the UserRepository
func NewUserRepository(db *gorm.DB) usersRepository {
	return usersRepository{storage: db}
}

func (repo usersRepository) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	err := repo.storage.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return handleErr(err)
	}

	return user, nil
}

func (repo usersRepository) FindUserById(ID int) model.User {
	var user model.User
	repo.storage.Where("id = ?", ID).First(&user)

	return user
}

func (repo usersRepository) StoreUser(user model.User) error {
	return repo.storage.Create(&user).Error
}

func handleErr(err error) (model.User, error) {
	if gorm.IsRecordNotFoundError(err) {
		return model.User{}, nil
	}

	return model.User{}, err
}
