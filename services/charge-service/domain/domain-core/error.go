package domain

import (
	"errors"
)

var (
	ErrInvalidChargeAmount  = errors.New("invalid charge amount")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrInvalidServiceID     = errors.New("invalid service ID")
	ErrAggregatorNotStarted = errors.New("charge aggregator not started")
)
