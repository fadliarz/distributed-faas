package domain

import (
	"fmt"
)

var (
	ErrInvalidBillingPeriod = fmt.Errorf("invalid billing period")
	ErrInvalidAmount        = fmt.Errorf("invalid billing amount")
	ErrInvalidUserID        = fmt.Errorf("invalid user ID")
	ErrInvalidBillingID     = fmt.Errorf("invalid billing ID")
	ErrEmptyChargeList      = fmt.Errorf("empty charge list for billing calculation")
)
