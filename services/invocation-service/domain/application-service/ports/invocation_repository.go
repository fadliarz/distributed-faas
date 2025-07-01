package ports

import (
	"context"

	"github.com/fadliarz/services/invocation-service/domain/domain-core"
)

type InvocationRepository interface {
	Save(ctx context.Context, invocation *domain.Invocation) error
}
