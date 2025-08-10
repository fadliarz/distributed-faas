package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
)

// Ports

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (valueobject.UserID, error)
}

// Interfaces

type UserDataMapper interface {
	CreateUserCommandToUser(command *CreateUserCommand) (*domain.User, error)
}

type UserApplicationService interface {
	PersistsUser(ctx context.Context, command *CreateUserCommand) (*domain.User, error)
}
