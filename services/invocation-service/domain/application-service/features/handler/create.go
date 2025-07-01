package handler

import (
	"context"

	"github.com/fadliarz/services/invocation-service/domain/application-service/features/command"
	"github.com/fadliarz/services/invocation-service/domain/domain-core"
	"github.com/fadliarz/services/invocation-service/domain/domain-core/core"
)

func (h *CommandHandler) CreateInvocation(ctx context.Context, cmd *command.CreateInvocationCommand) (domain.InvocationID, error) {
	// ToDo: authorize user using userID and functionID

	if err := validateCommand(cmd); err != nil {
		return "", err
	}

	invocationID, err := h.service.PersistInvocation(ctx, cmd)
	if err != nil {
		return "", err
	}

	return invocationID, nil
}

func validateCommand(cmd *command.CreateInvocationCommand) error {
	if cmd == nil {
		return core.NewInternalError("command cannot be nil", nil)
	}

	if cmd.FunctionID == "" {
		return core.NewValidationError("function ID cannot be empty", nil)
	}

	return nil
}
