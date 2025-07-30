package domain

import (
	"errors"
	"fmt"
)

var (
	ErrCheckpointAlreadyExists      = errors.New("checkpoint already exists")
	ErrCheckpointAlreadyReprocessed = errors.New("checkpoint already reprocessed")
)

func NewErrCheckpointAlreadyExists(err error) error {
	return fmt.Errorf("%w: %s", ErrCheckpointAlreadyExists, err.Error())
}

func NewErrCheckpointAlreadyReprocessed(err error) error {
	return fmt.Errorf("%w: %s", ErrCheckpointAlreadyReprocessed, err.Error())
}
