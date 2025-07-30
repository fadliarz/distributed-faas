package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvocationDataAccessMapperImpl struct{}

func NewInvocationDataAccessMapper() InvocationDataAccessMapper {
	return &InvocationDataAccessMapperImpl{}
}

func (m *InvocationDataAccessMapperImpl) Entity(invocation *domain.Invocation) (*InvocationEntity, error) {
	if invocation == nil {
		return nil, nil
	}

	invocationID, err := primitive.ObjectIDFromHex(invocation.InvocationID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid invocation ID %s format: %w", invocation.InvocationID.String(), err)
	}

	return &InvocationEntity{
		InvocationID:  invocationID,
		FunctionID:    invocation.FunctionID.String(),
		UserID:        invocation.UserID.String(),
		SourceCodeURL: invocation.SourceCodeURL.String(),
		OutputURL:     invocation.OutputURL.String(),
		Timestamp:     invocation.Timestamp.Int64(),
		IsRetry:       invocation.IsRetry.Bool(),
	}, nil
}

func (m *InvocationDataAccessMapperImpl) Domain(entity *InvocationEntity) *domain.Invocation {
	if entity == nil {
		return nil
	}

	return &domain.Invocation{
		InvocationID:  domain.NewInvocationID(entity.InvocationID.Hex()),
		FunctionID:    domain.NewFunctionID(entity.FunctionID),
		UserID:        domain.NewUserID(entity.UserID),
		SourceCodeURL: domain.NewSourceCodeURL(entity.SourceCodeURL),
		OutputURL:     domain.NewOutputURL(entity.OutputURL),
		Timestamp:     domain.NewTimestamp(entity.Timestamp),
		IsRetry:       domain.NewIsRetry(entity.IsRetry),
	}
}
