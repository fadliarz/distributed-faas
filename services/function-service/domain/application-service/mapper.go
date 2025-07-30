package application

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

type FunctionDataMapperImpl struct{}

func NewFunctionDataMapper() FunctionDataMapper {
	return &FunctionDataMapperImpl{}
}

func (m *FunctionDataMapperImpl) CreateFunctionCommandToFunction(cmd *CreateFunctionCommand) (*domain.Function, error) {
	return &domain.Function{
		UserID:        domain.NewUserID(cmd.UserID),
		SourceCodeURL: domain.NewSourceCodeURL(cmd.SourceCodeURL),
	}, nil
}
