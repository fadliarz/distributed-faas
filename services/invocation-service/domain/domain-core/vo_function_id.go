package domain

import "fmt"

type FunctionID string

func NewFunctionID(id string) (FunctionID, error) {
	if id == "" {
		return "", fmt.Errorf("function ID cannot be empty")
	}
	return FunctionID(id), nil
}

func (f FunctionID) String() string {
	return string(f)
}
