package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FunctionDataAccessMapperImpl struct{}

func NewFunctionDataAccessMapper() FunctionDataAccessMapper {
	return &FunctionDataAccessMapperImpl{}
}

func (m *FunctionDataAccessMapperImpl) Entity(function *domain.Function) (*FunctionEntity, error) {
	if function == nil {
		return nil, nil
	}

	functionID, err := primitive.ObjectIDFromHex(function.FunctionID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid function ID: %w", err)
	}

	return &FunctionEntity{
		FunctionID:    functionID,
		UserID:        function.UserID.String(),
		SourceCodeURL: function.SourceCodeURL.String(),
	}, nil
}

func (m *FunctionDataAccessMapperImpl) Domain(entity *FunctionEntity) *domain.Function {
	if entity == nil {
		return nil
	}

	return &domain.Function{
		UserID:        domain.NewUserID(entity.UserID),
		FunctionID:    domain.NewFunctionID(entity.FunctionID.Hex()),
		SourceCodeURL: domain.NewSourceCodeURL(entity.SourceCodeURL),
	}
}
