package repository

import (
	"github.com/fadliarz/services/invocation-service/domain/domain-core"
)

type InvocationRepository struct{}

func NewInvocationRepository() *InvocationRepository {
	return &InvocationRepository{}
}

func (r *InvocationRepository) Save(invocation *domain.Invocation) error {
	mapper := NewInvocationMapper()
	entity := mapper.Entity(invocation)

	mongo := NewInvocationMongoRepository()
	return mongo.Save(entity)
}
