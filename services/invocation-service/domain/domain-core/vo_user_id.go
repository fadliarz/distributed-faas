package domain

import "github.com/fadliarz/services/invocation-service/domain/domain-core/core"

type UserID string

func NewUserID(id string) (UserID, error) {
	if id == "" {
		return "", core.NewValidationError("user id cannot be empty", nil)
	}
	return UserID(id), nil
}

func (u UserID) String() string {
	return string(u)
}
