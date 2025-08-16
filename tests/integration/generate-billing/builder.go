package main

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewCronEntity(userID primitive.ObjectID, lastBilled int64) *CronEntity {
	return &CronEntity{
		UserID:     userID,
		LastBilled: lastBilled,
	}
}
