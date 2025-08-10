package domain

import "github.com/fadliarz/distributed-faas/common/valueobject"

type UserDomainService interface {
	ValidateAndInitiateUser(user *User, userID valueobject.UserID) error
}
