package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constructors

type InvocationApplicationService struct {
	mapper            InvocationDataMapper
	domainService     domain.InvocationDomainService
	repositoryManager *RepositoryManager
}

type RepositoryManager struct {
	Invocation InvocationRepository
	Function   FunctionRepository
}

func NewInvocationApplicationService(mapper InvocationDataMapper, domainSvc domain.InvocationDomainService, repositoryManager *RepositoryManager) *InvocationApplicationService {
	return &InvocationApplicationService{
		mapper:            mapper,
		domainService:     domainSvc,
		repositoryManager: repositoryManager,
	}
}

func NewInvocationApplicationServiceRepositoryManager(invocation InvocationRepository, function FunctionRepository) *RepositoryManager {
	return &RepositoryManager{
		Invocation: invocation,
		Function:   function,
	}
}

// Methods

func (s *InvocationApplicationService) PersistInvocation(ctx context.Context, cmd *CreateInvocationCommand) (*domain.Invocation, error) {
	function, err := s.repositoryManager.Function.FindByUserIDAndFunctionID(ctx, domain.UserID(cmd.UserID), domain.FunctionID(cmd.FunctionID))
	if err != nil {
		return nil, fmt.Errorf("failed to find function: %w", err)
	}

	if function == nil {
		return nil, domain.NewErrUserNotAuthorized(err)
	}

	invocation := s.mapper.CreateInvocationCommandToInvocation(cmd)

	if err := s.domainService.ValidateAndInitiateInvocation(invocation, primitive.NewObjectID().Hex(), function); err != nil {
		return nil, fmt.Errorf("failed to validate and initiate invocation: %w", err)
	}

	_, err = s.repositoryManager.Invocation.Save(ctx, invocation)
	if err != nil {
		return nil, fmt.Errorf("failed to save invocation: %w", err)
	}

	return invocation, nil
}

func (s *InvocationApplicationService) GetInvocation(ctx context.Context, query *GetInvocationQuery) (*domain.Invocation, error) {
	invocation, err := s.repositoryManager.Invocation.FindByID(ctx, domain.InvocationID(query.InvocationID))
	if err != nil {
		return nil, err
	}

	if invocation == nil || invocation.UserID.String() != query.UserID {
		return nil, domain.NewErrUserNotAuthorized(nil)
	}

	return invocation, nil
}
