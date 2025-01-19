package errors

import "errors"

type ErrInvalidStorageOperation interface {
	Error() string
	Operation() string
}

var ErrPostNotFound = errors.New("post not found")
