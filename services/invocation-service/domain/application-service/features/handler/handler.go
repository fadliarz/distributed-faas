package handler

import (
	"github.com/fadliarz/services/invocation-service/domain/application-service/service"
)

type CommandHandler struct {
	service *service.InvocationApplicationService
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{service: service.NewInvocationApplicationService()}
}
