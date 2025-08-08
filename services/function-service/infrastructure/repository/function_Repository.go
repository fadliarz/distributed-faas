package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
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
	functionEntity, err := r.mapper.Entity(function)
	if err != nil {
		return "", fmt.Errorf("failed to map function entity: %w", err)
	}

	functionID, err := r.repo.Save(ctx, functionEntity)

	return domain.FunctionID(functionID), err
}

func (r *FunctionRepositoryImpl) FindByUserIDAndFunctionID(ctx context.Context, userID domain.UserID, functionID domain.FunctionID) (*domain.Function, error) {
	functionObjectID, err := primitive.ObjectIDFromHex(functionID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid function ID format: %w", err)
	}

	entity, err := r.repo.FindByUserIDAndFunctionID(ctx, string(userID), functionObjectID)

	return r.mapper.Domain(entity), err
}

func (r *FunctionRepositoryImpl) UpdateSourceCodeURLByUserIDAndFunctionID(ctx context.Context, userID domain.UserID, functionID domain.FunctionID, sourceCodeURL domain.SourceCodeURL) error {
	primitiveFunctionID, err := primitive.ObjectIDFromHex(functionID.String())
	if err != nil {
		return fmt.Errorf("invalid function ID format: %w", err)
	}

	return r.repo.UpdateSourceCodeURLByUserIDAndFunctionID(ctx, userID.String(), primitiveFunctionID, sourceCodeURL.String())
}
