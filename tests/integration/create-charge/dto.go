package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChargeEvent struct {
	UserID           string `json:"user_id"`
	ServiceID        string `json:"service_id"`
	AggregatedAmount int64  `json:"aggregated_amount"`
}

type ChargeEntity struct {
	ChargeID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string             `bson:"user_id"`
	ServiceID         string             `bson:"service_id"`
	AccumulatedAmount int64              `bson:"accumulated_amount"`
}
