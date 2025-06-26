package handler

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/service"
)

type CommandHandler struct {
	service *service.FunctionApplicationService
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{service: service.NewFunctionApplicationService()}
}
