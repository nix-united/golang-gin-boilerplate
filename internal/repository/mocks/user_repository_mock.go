package mocks

import (
	"errors"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
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

func (u *UserRepositoryMock) FindUserByID(id int) model.User {
	//nolint:gosec // G115: integer overflow conversion int -> uint (gosec)
	if u.User.ID != uint(id) {
		return model.User{}
	}
	return *u.User
}

func (u *UserRepositoryMock) StoreUser(model.User) error {
	return nil
}

func NewUserRepositoryMock(user *model.User) *UserRepositoryMock {
	return &UserRepositoryMock{
		User: user,
	}
}
