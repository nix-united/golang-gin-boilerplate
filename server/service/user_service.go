package service

import (
	"basic_server/server/model"
	"basic_server/server/repository"
	"basic_server/server/request"
	"basic_server/server/utils"
)

// UserService provides a use case level for the user entity
type UserService interface {
	// Create takes a request with new user credentials and registers it.
	// An error will be returned if a user exists in the system, or
	// if an error occurs during interaction with the database.
	CreateUser(req request.RegisterRequest, en utils.Encryptor) error
}

type userService struct {
	userRepo repository.UsersRepository
}

// NewUserService returns an instance of the UserService
func NewUserService(ur repository.UsersRepository) UserService {
	return userService{
		userRepo: ur,
	}
}

func (srv userService) CreateUser(req request.RegisterRequest, en utils.Encryptor) error {
	var err error
	var user model.User

	user, err = srv.userRepo.FindUserByEmail(req.Email)

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

	err = srv.userRepo.StoreUser(model.User{
		Email:    req.Email,
		Password: encryptedPassword,
		FullName: req.FullName,
	})

	if err != nil {
		return err
	}

	return nil
}
