package service

import "github.com/nix-united/golang-gin-boilerplate/internal/errors"

type RestError struct {
	Status int   `json:"Status"`
	Error  error `json:"Error"`
}

type errUserAlreadyExists struct {
	message   string
	operation string
}

func NewErrUserAlreadyExists(msg, opName string) errors.ErrInvalidStorageOperation {
	return errUserAlreadyExists{
		message:   msg,
		operation: opName,
	}
}

func (e errUserAlreadyExists) Error() string {
	return e.message
}

func (e errUserAlreadyExists) Operation() string {
	return e.operation
}
