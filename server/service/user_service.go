package service

import "basic_server/server/model"

type NewUserService struct {}

func (service *NewUserService) CreateUser (email, password string) model.User {
	return model.User {
		Email: email,
		Password: password,
		FullName: "",
	}
}