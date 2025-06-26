package service

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/ports"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
	"github.com/fadliarz/distributed-faas/services/function-service/infrastructure/repository"
)

type FunctionApplicationService struct {
	mapper    *mapper
	domainSvc *domain.FunctionDomainService

	functionRepo ports.FunctionRepository
}

func NewFunctionApplicationService() *FunctionApplicationService {
	return &FunctionApplicationService{mapper: &mapper{}, domainSvc: domain.NewFunctionDomainService(), functionRepo: repository.NewFunctionRepository()}
}

func (s *FunctionApplicationService) PersistFunction(cmd *command.CreateFunctionCommand) (domain.FunctionID, error) {
	function, err := s.mapper.CreateFunctionCommandToFunction(cmd)
	if err != nil {
		return "", err
	}

	if err := s.domainSvc.ValidateAndInitiateFunction(function); err != nil {
		return "", err
	}

	if err := s.functionRepo.Save(function); err != nil {
		return "", err
	}

	return function.FunctionID, nil
}
