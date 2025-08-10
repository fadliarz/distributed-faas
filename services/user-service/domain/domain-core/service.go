package domain

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type UserDomainServiceImpl struct{}

func NewUserDomainService() UserDomainService {
	return &UserDomainServiceImpl{}
}

func (u *UserDomainServiceImpl) ValidateAndInitiateUser(user *User, userID valueobject.UserID) error {
	if user.UserID.String() == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	user.UserID = userID

	if user.Password.String() == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}
