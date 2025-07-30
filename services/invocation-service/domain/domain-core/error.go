package domain

import (
	"errors"
	"fmt"
)

var (
	ErrFunctionNotFound      = errors.New("function not found")
	ErrFunctionAlreadyExists = errors.New("function already exists")
	ErrUserNotAuthorized     = errors.New("user not authorized to perform this action")
	ErrInvocationNotFound    = errors.New("invocation not found")
)

func NewErrFunctionNotFound(err error) error {
	return fmt.Errorf("%w: %s", ErrFunctionNotFound, err.Error())
}

func NewErrFunctionAlreadyExists(err error) error {
	return fmt.Errorf("%w: %s", ErrFunctionAlreadyExists, err.Error())
}

func NewErrUserNotAuthorized(err error) error {
	return fmt.Errorf("%w: %s", ErrUserNotAuthorized, err.Error())
}

