package service

import (
	"basic_server/internal/model"
	"basic_server/repository"
	"basic_server/request"
	"basic_server/utils"
)

// UserServiceI provides a use case level for the user entity
type UserServiceI interface {
	// CreateUser Create takes a request with new user credentials and registers it.
	// An error will be returned if a user exists in the system, or
	// if an error occurs during interaction with the database.
	CreateUser(req request.RegisterRequest, en utils.Encryptor) error
}

type UserService struct {
	UserRepo repository.UserRepositoryI
}

// NewUserService returns an instance of the UserService
func NewUserService(ur repository.UserRepositoryI) UserServiceI {
	return &UserService{
		UserRepo: ur,
	}
}

func (srv UserService) CreateUser(req request.RegisterRequest, en utils.Encryptor) error {
	user, err := srv.UserRepo.FindUserByEmail(req.Email)

	if err != nil {
		return err
	}

	if user.ID != 0 {
		return NewErrUserAlreadyExists(
			"user already exist",
			"store a user",
		)
	}

	var encryptedPassword string

	encryptedPassword, err = en.Encrypt(req.Password)

	if err != nil {
		return err
	}

	err = srv.UserRepo.StoreUser(model.User{
		Email:    req.Email,
		Password: encryptedPassword,
		FullName: req.FullName,
	})

	if err != nil {
		return err
	}

	return nil
}
