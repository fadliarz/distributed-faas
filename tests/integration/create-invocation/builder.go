package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckpointEntity struct {
	CheckpointID  primitive.ObjectID `bson:"_id,omitempty"`
	FunctionID    string             `bson:"function_id"`
	UserID        string             `bson:"user_id"`
	SourceCodeURL string             `bson:"source_code_url"`
	Timestamp     int64              `bson:"timestamp"`
	Status        string             `bson:"status"`
	OutputURL     string             `bson:"output_url"`
}

func NewRandomCheckpointEntity() *CheckpointEntity {
	return &CheckpointEntity{
		CheckpointID:  primitive.NewObjectID(),
		FunctionID:    primitive.NewObjectID().Hex(),
		UserID:        primitive.NewObjectID().Hex(),
		SourceCodeURL: "user-id-123/main.js",
		Timestamp:     time.Now().Add(-5 * time.Minute).Unix(),
		Status:        "PENDING",
		OutputURL:     "",
	}
}
