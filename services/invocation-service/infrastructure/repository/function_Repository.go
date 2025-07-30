package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FunctionRepositoryImpl struct {
	mapper FunctionDataAccessMapper
	repo   *FunctionMongoRepository
}

func NewFunctionRepository(mapper FunctionDataAccessMapper, repo *FunctionMongoRepository) application.FunctionRepository {
	return &FunctionRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *FunctionRepositoryImpl) Save(ctx context.Context, function *domain.Function) (domain.FunctionID, error) {
	entity, err := r.mapper.Entity(function)
	if err != nil {
		return "", fmt.Errorf("failed to map function to entity: %w", err)
	}

	id, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", err
	}

	return domain.FunctionID(id), nil
}

func (r *FunctionRepositoryImpl) FindByUserIDAndFunctionID(ctx context.Context, userID domain.UserID, functionID domain.FunctionID) (*domain.Function, error) {
	functionObjectID, err := primitive.ObjectIDFromHex(functionID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid function ID format: %w", err)
	}

	entity, err := r.repo.FindByUserIDAndFunctionID(ctx, string(userID), functionObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find function by user ID and function ID: %w", err)
	}

	return r.mapper.Domain(entity), nil
}
