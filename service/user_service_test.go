package service

import (
	"errors"
	"testing"

	"basic_server/server/model"
	rmock "basic_server/server/repository/mocks"
	"basic_server/server/request"
	emock "basic_server/server/utils/mocks"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name            string
		arg             request.RegisterRequest
		encryptReceives string
		encryptReturns  struct {
			pass string
			err  error
		}
		findUserByEmailReceives string
		findUserByEmailReturns  struct {
			user model.User
			err  error
		}
		storeReceives model.User
		storeReturns  error
		wantErr       bool
	}{
		{
			"test successful creating a new user",
			request.RegisterRequest{
				Email:    "test@test.com",
				Password: "test pass",
				FullName: "test full name",
			},
			"test pass",
			struct {
				pass string
				err  error
			}{
				pass: "encrypted test pass",
				err:  nil,
			},
			"test@test.com",
			struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  nil,
			},
			model.User{
				Email:    "test@test.com",
				Password: "encrypted test pass",
				FullName: "test full name",
			},
			nil,
			false,
		},
		{
			"test returning an error if a user exists in the database",
			request.RegisterRequest{
				Email:    "test@test.com",
				Password: "",
				FullName: "",
			},
			"",
			struct {
				pass string
				err  error
			}{},
			"test@test.com",
			struct {
				user model.User
				err  error
			}{
				user: model.User{
					Model: gorm.Model{
						ID: 1,
					},
				},
				err: nil,
			},
			model.User{},
			nil,
			true,
		},
		{
			"test returning an error if it occurs during finding a user in the database",
			request.RegisterRequest{
				Email:    "test@test.com",
				Password: "",
				FullName: "",
			},
			"",
			struct {
				pass string
				err  error
			}{},
			"test@test.com",
			struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  errors.New("test"),
			},
			model.User{},
			nil,
			true,
		},
		{
			"test returning an error if it occurs during password encryption",
			request.RegisterRequest{
				Email:    "test@test.com",
				Password: "test pass",
				FullName: "test full name",
			},
			"test pass",
			struct {
				pass string
				err  error
			}{
				pass: "",
				err:  errors.New("test"),
			},
			"test@test.com",
			struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  nil,
			},
			model.User{},
			nil,
			true,
		},
		{
			"test returning an error if it occurs during saving a user to the database",
			request.RegisterRequest{
				Email:    "test@test.com",
				Password: "test pass",
				FullName: "test full name",
			},
			"test pass",
			struct {
				pass string
				err  error
			}{
				pass: "encrypted test pass",
				err:  nil,
			},
			"test@test.com",
			struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  nil,
			},
			model.User{
				Email:    "test@test.com",
				Password: "encrypted test pass",
				FullName: "test full name",
			},
			errors.New("test"),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor := &emock.Encryptor{}
			encryptor.
				On("Encrypt", tt.encryptReceives).
				Return(tt.encryptReturns.pass, tt.encryptReturns.err)

			userRepo := &rmock.UsersRepository{}
			userRepo.
				On("FindUserByEmail", tt.findUserByEmailReceives).
				Return(tt.findUserByEmailReturns.user, tt.findUserByEmailReturns.err).
				On("StoreUser", tt.storeReceives).
				Return(tt.storeReturns)

			err := NewUserService(userRepo).CreateUser(tt.arg, encryptor)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
