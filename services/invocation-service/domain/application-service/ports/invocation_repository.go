package ports

import "github.com/fadliarz/services/invocation-service/domain/domain-core"

type InvocationRepository interface {
	Save(invocation *domain.Invocation) error
}
