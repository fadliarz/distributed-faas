package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constructors

type UserApplicationServiceImpl struct {
	mapper            UserDataMapper
	domainService     domain.UserDomainService
	repositoryManager *UserApplicationServiceRepositoryManager
}

type UserApplicationServiceRepositoryManager struct {
	User UserRepository
}

func NewUserApplicationService(mapper UserDataMapper, domainService domain.UserDomainService, repositoryManager *UserApplicationServiceRepositoryManager) UserApplicationService {
	return &UserApplicationServiceImpl{
		mapper:            mapper,
		domainService:     domainService,
		repositoryManager: repositoryManager,
	}
}

func NewUserApplicationServiceRepositoryManager(user UserRepository) *UserApplicationServiceRepositoryManager {
	return &UserApplicationServiceRepositoryManager{
		User: user,
	}
}

// Methods

func (s *UserApplicationServiceImpl) PersistsUser(ctx context.Context, command *CreateUserCommand) (*domain.User, error) {
	user, err := s.mapper.CreateUserCommandToUser(command)
	if err != nil {
		return nil, fmt.Errorf("failed to map command to user: %w", err)
	}

	err = s.domainService.ValidateAndInitiateUser(user, valueobject.NewUserID(primitive.NewObjectID().Hex()))
	if err != nil {
		return nil, fmt.Errorf("failed to validate and initiate user: %w", err)
	}

	_, err = s.repositoryManager.User.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}
