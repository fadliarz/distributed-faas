package domain

import (
	"errors"

	"github.com/google/uuid"
)

type FunctionDomainService struct{}

func NewFunctionDomainService() *FunctionDomainService {
	return &FunctionDomainService{}
}

func (s *FunctionDomainService) ValidateAndInitiateFunction(function *Function) error {
	defaultErr := errors.New("failed validating and initializing function")

	if function.FunctionID.String() != "" {
		return defaultErr
	}

	functionId, err := NewFunctionID(uuid.NewString())
	if err != nil {
		return defaultErr
	}
	function.FunctionID = functionId

	return nil
}
