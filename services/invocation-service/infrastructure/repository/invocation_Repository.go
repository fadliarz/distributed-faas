package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvocationRepository struct {
	mapper InvocationDataAccessMapper
	repo   *InvocationMongoRepository
}

func NewInvocationRepository(mapper InvocationDataAccessMapper, repo *InvocationMongoRepository) application.InvocationRepository {
	return &InvocationRepository{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *InvocationRepository) Save(ctx context.Context, invocation *domain.Invocation) (domain.InvocationID, error) {
	entity, err := r.mapper.Entity(invocation)
	if err != nil {
		return "", fmt.Errorf("failed to map invocation to entity: %w", err)
	}

	id, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", err
	}

	return domain.InvocationID(id), nil
}

func (r *InvocationRepository) FindByID(ctx context.Context, invocationID domain.InvocationID) (*domain.Invocation, error) {
	primitiveInvocationID, err := primitive.ObjectIDFromHex(invocationID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse invocation ID: %w", err)
	}

	entity, err := r.repo.FindByID(ctx, primitiveInvocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find invocation by ID: %w", err)
	}

	invocation := r.mapper.Domain(entity)

	return invocation, nil
}
