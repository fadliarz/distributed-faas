package handler

import (
	"github.com/fadliarz/services/invocation-service/domain/application-service/features/command"
	"github.com/fadliarz/services/invocation-service/domain/domain-core"
	"github.com/fadliarz/services/invocation-service/domain/domain-core/core"
)

func (h *CommandHandler) CreateInvocation(cmd *command.CreateInvocationCommand) (domain.InvocationID, error) {
	if err := validateCommand(cmd); err != nil {
		return "", err
	}

	invocationID, err := h.service.PersistInvocation(cmd)
	if err != nil {
		return "", err
	}

	return invocationID, nil
}

func validateCommand(cmd *command.CreateInvocationCommand) error {
	if cmd == nil {
		return core.NewInternalError("command cannot be nil", nil)
	}

	if cmd.UserID == "" {
		return core.NewValidationError("user ID cannot be empty", nil)
	}

	if cmd.FunctionID == "" {
		return core.NewValidationError("function ID cannot be empty", nil)
	}
	return nil
}
