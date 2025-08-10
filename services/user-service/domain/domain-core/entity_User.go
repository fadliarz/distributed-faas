package domain

import "github.com/fadliarz/distributed-faas/common/valueobject"

type User struct {
	UserID   valueobject.UserID
	Password valueobject.Password
}
