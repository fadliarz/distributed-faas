package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/domain-core"
)

// Ports

type MachineRepository interface {
	FindManyByStatus(ctx context.Context, status domain.Status) ([]domain.Machine, error)
}

type InvocationConsumer interface {
	common.Consumer[*InvocationCreatedEvent]
}
