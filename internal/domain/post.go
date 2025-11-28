package domain

import (
	"errors"
	"fmt"
)

const maxLimit = 100

type PostFilters struct {
	Offset int64
	Limit  int64
}

func (f PostFilters) Validate() error {
	var err error
	if f.Limit < 0 {
		err = errors.Join(err, errors.New("limit should be positive number"))
	}

	if f.Limit > maxLimit {
		err = errors.Join(err, fmt.Errorf("limit should be not more than %d", maxLimit))
	}

	if f.Offset < 0 {
		err = errors.Join(err, errors.New("offset should be positive number"))
	}

	return err
}
