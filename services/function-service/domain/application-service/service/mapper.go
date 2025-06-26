package service

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

type mapper struct{}

func (m *mapper) CreateFunctionCommandToFunction(cmd *command.CreateFunctionCommand) (*domain.Function, error) {
	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	sourceCodeURL, err := domain.NewSourceCodeURL(cmd.SourceCodeURL)
	if err != nil {
		return nil, err
	}

	return &domain.Function{
		UserID:        userID,
		SourceCodeURL: sourceCodeURL,
	}, nil
}
