package application

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type ChargeEvent struct {
	UserID    valueobject.UserID    `json:"user_id"`
	ServiceID valueobject.ServiceID `json:"service_id"`
	Amount    valueobject.Amount    `json:"amount"`
}
