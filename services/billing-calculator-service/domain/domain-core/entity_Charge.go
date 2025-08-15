package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type Charge struct {
	ChargeID          valueobject.ChargeID
	UserID            valueobject.UserID
	ServiceID         valueobject.ServiceID
	AccumulatedAmount valueobject.Amount
}
