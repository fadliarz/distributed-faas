package domain

import (
	"github.com/fadliarz/services/invocation-service/domain/domain-core/core"
	"github.com/google/uuid"
)

type InvocationID string

func NewInvocationID(id string) (InvocationID, error) {
	if id == "" {
		return "", core.NewValidationError("invocation id cannot be empty", nil)
	}
	return InvocationID(id), nil
}

func GenerateInvocationID() (InvocationID, error) {
	id := uuid.New().String()
	return NewInvocationID(id)
}

func (i InvocationID) String() string {
	return string(i)
}
