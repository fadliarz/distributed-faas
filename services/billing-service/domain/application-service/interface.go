package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-service/domain/domain-core"
)

// Ports

type BillingRepository interface {
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*domain.Billing, error)
}

// Interfaces

type BillingApplicationService interface {
	GetBilling(ctx context.Context, query *GetBillingQuery) (*domain.Billing, error)
}

// Commands and Queries

type GetBillingQuery struct {
	UserID string
}
