package application

import (
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
)

type InvocationDataMapperImpl struct{}

func NewInvocationDataMapper() InvocationDataMapper {
	return &InvocationDataMapperImpl{}
}

func (m *InvocationDataMapperImpl) CreateInvocationCommandToInvocation(cmd *CreateInvocationCommand) *domain.Invocation {
	return &domain.Invocation{
		FunctionID: domain.NewFunctionID(cmd.FunctionID),
		UserID:     domain.NewUserID(cmd.UserID),
	}
}
