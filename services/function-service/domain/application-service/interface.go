package application

import (
	"context"
	"time"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

// Ports

type FunctionRepository interface {
	Save(ctx context.Context, function *domain.Function) (domain.FunctionID, error)
	FindByUserIDAndFunctionID(ctx context.Context, userID domain.UserID, functionID domain.FunctionID) (*domain.Function, error)
	UpdateSourceCodeURLByUserIDAndFunctionID(ctx context.Context, userID domain.UserID, functionID domain.FunctionID, sourceCodeURL domain.SourceCodeURL) error
}

// Interfaces

type FunctionDataMapper interface {
	CreateFunctionCommandToFunction(cmd *CreateFunctionCommand) (*domain.Function, error)
}

type FunctionApplicationService interface {
	PersistFunction(ctx context.Context, command *CreateFunctionCommand) (*domain.Function, error)
	GetFunctionUploadPresignedURL(ctx context.Context, query *GetFunctionUploadPresignedURLQuery) (string, error)
	UpdateFunctionSourceCodeURL(ctx context.Context, command *UpdateFunctionSourceCodeURLCommand) error
}

type InputStorage interface {
	GetFunctionUploadPresignedURL(ctx context.Context, userID domain.UserID, functionID domain.FunctionID, language domain.Language, expiration time.Duration) (string, error)
}
