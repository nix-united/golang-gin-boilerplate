package domain

import (
	"errors"
	"fmt"
)

const maxLimit = 100

type CreatePostRequest struct {
	UserID  uint
	Title   string
	Content string
}

type UpdatePostRequest struct {
	PostID  uint
	UserID  uint
	Title   string
	Content string
}

type PostFilters struct {
	UserID uint
	Title  string
	Offset int
	Limit  int
}

func (f PostFilters) Validate() error {
	var err error
	if f.Limit < 1 {
		err = errors.Join(err, errors.New("limit should be at least 1"))
	}

	if f.Limit > maxLimit {
		err = errors.Join(err, fmt.Errorf("limit should be not more than %d", maxLimit))
	}

	if f.Offset < 0 {
		err = errors.Join(err, errors.New("offset should be positive number"))
	}

	return err
}
