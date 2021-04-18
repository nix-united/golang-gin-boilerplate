package mocks

import (
	"basic_server/model"
	"basic_server/repository"
	"errors"
)

type UserRepositoryMock struct {
	User *model.User
}

func (u *UserRepositoryMock) FindUserByEmail(email string) (model.User, error) {
	if u.User.Email != email {
		return model.User{}, errors.New("not found")
	}
	return *u.User, nil
}

func (u *UserRepositoryMock) FindUserByID(ID int) model.User {
	if u.User.ID != uint(ID) {
		return model.User{}
	}
	return *u.User
}

func (u *UserRepositoryMock) StoreUser(*model.User) error {
	return nil
}

func NewUserRepositoryMock(user *model.User) repository.UserRepositoryI {
	return &UserRepositoryMock{
		User: user,
	}
}
