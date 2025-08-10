package domain

import "github.com/fadliarz/distributed-faas/common/valueobject"

type Cron struct {
	UserID      valueobject.UserID
	LastBilling valueobject.LastBilled
}
