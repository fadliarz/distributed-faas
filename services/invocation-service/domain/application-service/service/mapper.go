package service

import (
	"fmt"

	"github.com/fadliarz/services/invocation-service/domain/application-service/features/command"
	"github.com/fadliarz/services/invocation-service/domain/domain-core"
)

type mapper struct{}

func (m *mapper) CreateInvocationCommandToInvocation(cmd *command.CreateInvocationCommand) (*domain.Invocation, error) {
	functionID, err := domain.NewFunctionID(cmd.FunctionID)
	if err != nil {
		return nil, fmt.Errorf("invalid function ID: %w", err)
	}

	return &domain.Invocation{
		FunctionID: functionID,
	}, nil
}
