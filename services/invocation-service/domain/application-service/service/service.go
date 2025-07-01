package service

import (
	"context"
	"fmt"

	"github.com/fadliarz/services/invocation-service/domain/application-service/features/command"
	"github.com/fadliarz/services/invocation-service/domain/application-service/ports"
	"github.com/fadliarz/services/invocation-service/domain/domain-core"
	"github.com/fadliarz/services/invocation-service/infrastructure/repository"
)

type InvocationApplicationService struct {
	mapper         *mapper
	domainSvc      *domain.InvocationDomainService
	invocationRepo ports.InvocationRepository
}

func NewInvocationApplicationService() *InvocationApplicationService {
	return &InvocationApplicationService{mapper: &mapper{}, domainSvc: domain.NewInvocationDomainService(), invocationRepo: repository.NewInvocationRepository()}
}

func (s *InvocationApplicationService) PersistInvocation(ctx context.Context, cmd *command.CreateInvocationCommand) (domain.InvocationID, error) {
	invocation, err := s.mapper.CreateInvocationCommandToInvocation(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to map command to invocation: %w", err)
	}

	if err = s.domainSvc.ValidateAndInitiateInvocation(invocation); err != nil {
		return "", fmt.Errorf("failed to validate and initiate invocation: %w", err)
	}

	if err = s.invocationRepo.Save(ctx, invocation); err != nil {
		return "", err
	}

	return invocation.InvocationID, nil
}
