package handler

import "github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/service"

type CommandHandler struct {
	service *service.RegistrarApplicationService
}

func NewCommandHandler(service *service.RegistrarApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}
