package domain

import "github.com/fadliarz/services/invocation-service/domain/domain-core/core"

type InvocationDomainService struct{}

func NewInvocationDomainService() *InvocationDomainService {
	return &InvocationDomainService{}
}

func (s *InvocationDomainService) ValidateAndInitiateInvocation(invocation *Invocation) error {
	if invocation.FunctionID == "" {
		return core.NewValidationError("function ID is required", nil)
	}
	
	invocationID, err := GenerateInvocationID()
	if err != nil {
		return core.NewInternalError("failed to generate invocation ID", err)
	}
	invocation.InvocationID = invocationID
	return nil
}
