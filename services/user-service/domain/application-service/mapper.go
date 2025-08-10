package application

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
)

type UserDataMapperImpl struct{}

func NewUserDataMapper() UserDataMapper {
	return &UserDataMapperImpl{}
}

func (m *UserDataMapperImpl) CreateUserCommandToUser(command *CreateUserCommand) (*domain.User, error) {
	return &domain.User{
		Password: valueobject.Password(command.Password),
	}, nil
}
