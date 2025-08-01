package domain

import (
	"errors"
	"fmt"
)

var (
	ErrMachineNotFound = errors.New("machine not found")
)

func NewErrMachineNotFound(err error) error {
	return fmt.Errorf("%w: %w", ErrMachineNotFound, err)
}
