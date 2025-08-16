package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/billing-cron-service/domain/domain-core"
)

// Ports

type CronRepository interface {
	UpdateLastBilled(ctx context.Context, timestampPair domain.TimestampPair) error
}
