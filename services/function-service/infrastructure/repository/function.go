package repository

import (
	"errors"

	"github.com/fadliarz/services/function-service/domain/domain-core"
)

type FunctionRepositoryImpl struct {
	repo   *FunctionMongoRepository
	mapper *FunctionMapper
}

func NewFunctionRepository() *FunctionRepositoryImpl {
	return &FunctionRepositoryImpl{repo: NewFunctionMongoRepository(), mapper: NewFunctionMapper()}
}

func (r *FunctionRepositoryImpl) Save(function *domain.Function) error {
	defaultErr := errors.New("")

	functionEntity := r.mapper.Entity(function)

	err := r.repo.Save(functionEntity)
	if err != nil {
		return defaultErr
	}

	return nil
}
