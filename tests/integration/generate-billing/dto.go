package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type CronEntity struct {
	UserID     primitive.ObjectID `bson:"_id"`
	LastBilled int64              `bson:"last_billed"`
}

type CronEvent struct {
	UserID     string `json:"_id"`
	LastBilled int64  `json:"last_billed"`
}
