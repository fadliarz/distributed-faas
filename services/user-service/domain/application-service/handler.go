package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
)

// Constructor

type UserCommandHandler struct {
	service UserApplicationService
}

func NewUserCommandHandler(service UserApplicationService) *UserCommandHandler {
	return &UserCommandHandler{
		service: service,
	}
}

// Methods

func (h *UserCommandHandler) CreateUser(ctx context.Context, command *CreateUserCommand) (*domain.User, error) {
	return h.service.PersistsUser(ctx, command)
}
