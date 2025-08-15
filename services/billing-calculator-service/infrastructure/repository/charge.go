package repository

import (
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type ChargeEntity struct {
	ChargeID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string             `bson:"user_id"`
	ServiceID         string             `bson:"service_id"`
	AccumulatedAmount int64              `bson:"accumulated_amount"`
}

// Interfaces

type ChargeDataAccessMapper interface {
	Entity(charge *domain.Charge) (*ChargeEntity, error)
	Domain(entity *ChargeEntity) *domain.Charge
}
