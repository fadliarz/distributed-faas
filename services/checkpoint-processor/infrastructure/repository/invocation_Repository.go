package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvocationRepository struct {
	repo *InvocationMongoRepository
}

func NewInvocationRepository(repo *InvocationMongoRepository) application.InvocationRepository {
	return &InvocationRepository{
		repo: repo,
	}
}

func (r *InvocationRepository) UpdateOutputURLIfNotSet(ctx context.Context, invocationID domain.InvocationID, outputURL string) error {
	invocationIDPrimitive, err := primitive.ObjectIDFromHex(invocationID.String())
	if err != nil {
		return fmt.Errorf("invalid invocation ID: %w", err)
	}

	err = r.repo.UpdateOutputURLIfNotSet(ctx, invocationIDPrimitive, outputURL)
	if err != nil {
		return fmt.Errorf("failed to update output URL: %w", err)
	}

	return nil
}
