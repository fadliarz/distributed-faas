package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CronEntity struct {
	UserID      string `bson:"_id,omitempty"`
	LastBilling int64  `bson:"last_billing"`
}

func NewRandomCronEntity() *CronEntity {
	return &CronEntity{
		UserID:      primitive.NewObjectID().Hex(),
		LastBilling: time.Now().Unix(),
	}
}
