package domain

import (
	"fmt"
)

type FunctionDomainServiceImpl struct{}

func NewFunctionDomainService() FunctionDomainService {
	return &FunctionDomainServiceImpl{}
}

func (s *FunctionDomainServiceImpl) ValidateAndInitiateFunction(function *Function, functionID string) error {
	if function == nil {
		return fmt.Errorf("function cannot be nil")
	}

	if function.FunctionID != "" {
		return fmt.Errorf("function already has a FunctionID: %s", function.FunctionID)
	}

	if function.UserID == "" {
		return fmt.Errorf("function must have a UserID")
	}

	if function.SourceCodeURL == "" {
		return fmt.Errorf("function must have a SourceCodeURL")
	}

	function.FunctionID = NewFunctionID(functionID)

	return nil
}
