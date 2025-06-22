package repository

import (
	"errors"

	"github.com/fadliarz/services/function-service/domain/domain-core"
)

type FunctionMapper struct{}

func NewFunctionMapper() *FunctionMapper {
	return &FunctionMapper{}
}

func (m *FunctionMapper) Entity(function *domain.Function) *FunctionEntity {
	return &FunctionEntity{
		UserID:        function.UserID.String(),
		FunctionID:    function.FunctionID.String(),
		SourceCodeURL: function.FunctionID.String(),
	}
}

func (m *FunctionMapper) Domain(function FunctionEntity) (*domain.Function, error) {
	defaultErr := errors.New("")

	userID, err := domain.NewUserID(function.UserID)
	if err != nil {
		return nil, defaultErr
	}

	functionID, err := domain.NewFunctionID(function.FunctionID)
	if err != nil {
		return nil, defaultErr
	}

	sourceCodeURL, err := domain.NewSourceCodeURL(function.SourceCodeURL)
	if err != nil {
		return nil, defaultErr
	}

	return &domain.Function{
		UserID:        userID,
		FunctionID:    functionID,
		SourceCodeURL: sourceCodeURL,
	}, nil
}
