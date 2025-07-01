package domain

import (
	"fmt"
)

type InvocationDomainService struct{}

func NewInvocationDomainService() *InvocationDomainService {
	return &InvocationDomainService{}
}

func (s *InvocationDomainService) ValidateAndInitiateInvocation(invocation *Invocation) error {
	if invocation.FunctionID == "" {
		return fmt.Errorf("function ID cannot be empty")
	}

	invocationID, err := GenerateInvocationID()
	if err != nil {
		return fmt.Errorf("failed to generate invocation ID: %w", err)
	}

	invocation.InvocationID = invocationID

	return nil
}
