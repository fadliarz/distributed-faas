package domain

import "github.com/fadliarz/services/invocation-service/domain/domain-core/core"

type FunctionID string

func NewFunctionID(id string) (FunctionID, error) {
	if id == "" {
		return "", core.NewValidationError("function id cannot be empty", nil)
	}
	return FunctionID(id), nil
}

func (f FunctionID) String() string {
	return string(f)
}
