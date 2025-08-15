package repository

import (
	"github.com/fadliarz/distributed-faas/services/billing-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type BillingEntity struct {
	BillingID  primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	LastBilled int64              `bson:"last_billed"`
	Amount     int64              `bson:"amount"`
}

// Interfaces

type BillingDataAccessMapper interface {
	Entity(billing *domain.Billing) (*BillingEntity, error)
	Domain(entity *BillingEntity) *domain.Billing
}
