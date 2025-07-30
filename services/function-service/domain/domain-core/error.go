package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotAuthorized = errors.New("user not authorized")
	ErrFunctionNotFound  = errors.New("function not found")
)

func NewErrUserNotAuthorized(err error) error {
	return fmt.Errorf("%w: %s", ErrUserNotAuthorized, err.Error())
}

func NewErrFunctionNotFound(err error) error {
	return fmt.Errorf("%w: %s", ErrFunctionNotFound, err.Error())
}