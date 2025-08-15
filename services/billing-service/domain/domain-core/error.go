package domain

import (
	"errors"
	"fmt"
)

var (
	ErrBillingNotFound = errors.New("billing not found for the given user ID")
)

func NewErrBillingNotFound(err error) error {
	if err == nil {
		return ErrBillingNotFound
	}

	return fmt.Errorf("%w: %s", ErrBillingNotFound, err.Error())
}
