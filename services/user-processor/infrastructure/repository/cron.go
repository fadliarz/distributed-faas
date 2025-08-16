package repository

import (
	domain "github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type CronEntity struct {
	UserID      primitive.ObjectID `bson:"_id,omitempty"`
	LastBilling int64              `bson:"last_billing"`
}

// Mapper

type CronDataAccessMapper interface {
	Entity(cron *domain.Cron) (*CronEntity, error)
	Domain(entity *CronEntity) *domain.Cron
}
