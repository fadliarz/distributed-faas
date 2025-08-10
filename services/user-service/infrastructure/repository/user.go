package repository

import (
	"github.com/fadliarz/distributed-faas/services/user-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type UserEntity struct {
	UserID   primitive.ObjectID `bson:"_id,omitempty"`
	Password string             `bson:"password"`
}

// Interfaces

type UserDataAccessMapper interface {
	Entity(user *domain.User) (*UserEntity, error)
	Domain(user *UserEntity) *domain.User
}
