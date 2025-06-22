package service

import (
	"errors"

	"github.com/fadliarz/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/services/function-service/domain/application-service/ports"
	"github.com/fadliarz/services/function-service/domain/domain-core"
	"github.com/fadliarz/services/function-service/infrastructure/repository"
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
	defaultErr := errors.New("")

	function, err := s.mapper.CreateFunctionCommandToFunction(cmd)
	if err != nil {
		return "", defaultErr
	}

	s.domainSvc.ValidateAndInitiateFunction(function)

	s.functionRepo.Save(function)

	return domain.NewFunctionID("uuid")
}
