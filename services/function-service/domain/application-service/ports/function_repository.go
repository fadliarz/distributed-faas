package ports

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

type FunctionRepository interface {
	Save(function *domain.Function) error
}
