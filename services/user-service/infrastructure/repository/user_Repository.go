package repository

import (
	"context"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/user-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
)

type UserRepositoryImpl struct {
	mapper UserDataAccessMapper
	repo   *UserMongoRepository
}

func NewUserRepository(mapper UserDataAccessMapper, repo *UserMongoRepository) application.UserRepository {
	return &UserRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *domain.User) (valueobject.UserID, error) {
	entity, err := r.mapper.Entity(user)
	if err != nil {
		return "", err
	}

	userID, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", err
	}

	return valueobject.NewUserID(userID), nil
}
