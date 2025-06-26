package repository

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

type FunctionRepositoryImpl struct {
	repo   *FunctionMongoRepository
	mapper *FunctionMapper
}

func NewFunctionRepository() *FunctionRepositoryImpl {
	return &FunctionRepositoryImpl{repo: NewFunctionMongoRepository(), mapper: NewFunctionMapper()}
}

func (r *FunctionRepositoryImpl) Save(function *domain.Function) error {
	functionEntity := r.mapper.Entity(function)

	err := r.repo.Save(functionEntity)

	return err
}
