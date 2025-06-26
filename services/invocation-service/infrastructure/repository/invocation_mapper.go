package repository

import (
	"github.com/fadliarz/services/invocation-service/domain/domain-core"
)

type InvocationMapper struct{}

func NewInvocationMapper() *InvocationMapper {
	return &InvocationMapper{}
}

func (m *InvocationMapper) Entity(invocation *domain.Invocation) *InvocationEntity {
	return &InvocationEntity{
		InvocationID: invocation.InvocationID.String(),
		FunctionID:   invocation.FunctionID.String(),
	}
}

func (m *InvocationMapper) Domain(entity InvocationEntity) (*domain.Invocation, error) {
	invocationID, err := domain.NewInvocationID(entity.InvocationID)
	if err != nil {
		return nil, err
	}

	functionID, err := domain.NewFunctionID(entity.FunctionID)
	if err != nil {
		return nil, err
	}

	return &domain.Invocation{
		InvocationID: invocationID,
		FunctionID:   functionID,
	}, nil
}
