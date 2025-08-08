package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constructors

type FunctionApplicationServiceImpl struct {
	mapper            FunctionDataMapper
	domainService     domain.FunctionDomainService
	repositoryManager *FunctionApplicationServiceRepositoryManager
	storage           InputStorage
}

type FunctionApplicationServiceRepositoryManager struct {
	Function FunctionRepository
}

func NewFunctionApplicationService(mapper FunctionDataMapper, domainService domain.FunctionDomainService, repositoryManager *FunctionApplicationServiceRepositoryManager) FunctionApplicationService {
	return &FunctionApplicationServiceImpl{
		mapper:            mapper,
		domainService:     domainService,
		repositoryManager: repositoryManager,
	}
}

func NewFunctionApplicationServiceRepositoryManager(function FunctionRepository) *FunctionApplicationServiceRepositoryManager {
	return &FunctionApplicationServiceRepositoryManager{
		Function: function,
	}
}

// Methods

func (s *FunctionApplicationServiceImpl) PersistFunction(ctx context.Context, command *CreateFunctionCommand) (*domain.Function, error) {
	function, err := s.mapper.CreateFunctionCommandToFunction(command)
	if err != nil {
		return nil, fmt.Errorf("failed to map command to function: %w", err)
	}

	err = s.domainService.ValidateAndInitiateFunction(function, primitive.NewObjectID().Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to validate and initiate function: %w", err)
	}

	_, err = s.repositoryManager.Function.Save(ctx, function)
	if err != nil {
		return nil, fmt.Errorf("failed to save function: %w", err)
	}

	return function, nil
}

func (s *FunctionApplicationServiceImpl) GetFunctionUploadPresignedURL(ctx context.Context, query *GetFunctionUploadPresignedURLQuery) (string, error) {
	function, err := s.repositoryManager.Function.FindByUserIDAndFunctionID(ctx, domain.NewUserID(query.UserID), domain.NewFunctionID(query.FunctionID))
	if function == nil {
		return "", domain.NewErrUserNotAuthorized(err)
	}

	if err != nil {
		return "", fmt.Errorf("failed to find function by user ID and function ID: %w", err)
	}

	url, err := s.storage.GetFunctionUploadPresignedURL(ctx, domain.NewUserID(query.UserID), domain.NewFunctionID(query.FunctionID), domain.NewLanguage(query.Language), 1*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL: %w", err)
	}

	return url, nil
}

func (s *FunctionApplicationServiceImpl) UpdateFunctionSourceCodeURL(ctx context.Context, command *UpdateFunctionSourceCodeURLCommand) error {
	err := s.repositoryManager.Function.UpdateSourceCodeURLByUserIDAndFunctionID(ctx, domain.NewUserID(command.UserID), domain.NewFunctionID(command.FunctionID), domain.NewSourceCodeURL(command.SourceCodeURL))

	if err != nil && errors.Is(err, domain.ErrFunctionNotFound) {
		return domain.NewErrUserNotAuthorized(err)
	}

	return err
}
