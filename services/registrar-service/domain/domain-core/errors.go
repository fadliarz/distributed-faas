package domain

import "errors"

var (
	ErrUserNotAuthorized = errors.New("you're not authorized to perform this action")
)
