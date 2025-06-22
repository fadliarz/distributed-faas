package handler

import (
	"github.com/fadliarz/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/services/function-service/domain/domain-core"
)

func (h *CommandHandler) CreateFunction(cmd *command.CreateFunctionCommand) (domain.FunctionID, error) {
	h.service.PersistFunction(cmd)

	return "", nil
}
