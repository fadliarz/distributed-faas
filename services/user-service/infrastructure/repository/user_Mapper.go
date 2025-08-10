package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDataAccessMapperImpl struct{}

func NewUserDataAccessMapper() UserDataAccessMapper {
	return &UserDataAccessMapperImpl{}
}

func (m *UserDataAccessMapperImpl) Entity(user *domain.User) (*UserEntity, error) {
	if user == nil {
		return nil, nil
	}

	primitiveUserID, err := primitive.ObjectIDFromHex(user.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to convert userID to ObjectID: %w", err)
	}

	return &UserEntity{
		UserID:   primitiveUserID,
		Password: user.Password.String(),
	}, nil
}

func (m *UserDataAccessMapperImpl) Domain(entity *UserEntity) *domain.User {
	if entity == nil {
		return nil
	}

	return &domain.User{
		UserID:   valueobject.NewUserID(entity.UserID.Hex()),
		Password: valueobject.NewPassword(entity.Password),
	}
}
