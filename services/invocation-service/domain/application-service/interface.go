package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
)

// Ports

type FunctionRepository interface {
	Save(ctx context.Context, function *domain.Function) (domain.FunctionID, error)
	FindByUserIDAndFunctionID(ctx context.Context, userID domain.UserID, functionID domain.FunctionID) (*domain.Function, error)
}

type InvocationRepository interface {
	Save(ctx context.Context, invocation *domain.Invocation) (domain.InvocationID, error)
	FindByID(ctx context.Context, invocationID domain.InvocationID) (*domain.Invocation, error)
}

// Interfaces

type InvocationDataMapper interface {
	CreateInvocationCommandToInvocation(cmd *CreateInvocationCommand) *domain.Invocation
}
