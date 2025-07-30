package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvocationNotFound = errors.New("invocation not found")
)

func NewErrInvocationNotFound(err error) error {
	return fmt.Errorf("%w: %s", ErrInvocationNotFound, err.Error())
}
