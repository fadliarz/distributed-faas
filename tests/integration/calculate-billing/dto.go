package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChargeEntity struct {
	ChargeID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string             `bson:"user_id"`
	ServiceID         string             `bson:"service_id"`
	Timestamp         int64              `bson:"timestamp"`
	AccumulatedAmount int64              `bson:"accumulated_amount"`
}

type BillingEntity struct {
	BillingID primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	Amount    int64              `bson:"amount"`
}

type CronEvent struct {
	UserID     string `json:"_id"`
	LastBilled int64  `json:"last_billed"`
}
