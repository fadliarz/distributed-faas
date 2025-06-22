package service

import (
	"errors"

	"github.com/fadliarz/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/services/function-service/domain/domain-core"
)

type mapper struct{}

func (m *mapper) CreateFunctionCommandToFunction(cmd *command.CreateFunctionCommand) (*domain.Function, error) {
	defaultErr := errors.New("")

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		return nil, defaultErr
	}

	sourceCodeURL, err := domain.NewSourceCodeURL(cmd.SourceCodeURL)
	if err != nil {
		return nil, defaultErr
	}

	return &domain.Function{
		UserID:        userID,
		SourceCodeURL: sourceCodeURL,
	}, nil
}
