package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
)

// Ports

type CronRepository interface {
	Save(ctx context.Context, cron *domain.Cron) (domain.CronID, error)
}

type UserConsumer interface {
	common.Consumer[*UserEvent]
}

// Interfaces

type UserProcessorDataMapper interface {
	UserEventToCron(event *UserEvent) *domain.Cron
}
