package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type Charge struct {
	UserID    valueobject.UserID
	ServiceID valueobject.ServiceID
	Amount    valueobject.Amount
	Timestamp valueobject.Timestamp
}

type ChargeAggregate struct {
	UserID           valueobject.UserID
	ServiceID        valueobject.ServiceID
	AggregatedAmount valueobject.Amount
}

func NewChargeAggregate(userID valueobject.UserID, serviceID valueobject.ServiceID) *ChargeAggregate {
	return &ChargeAggregate{
		UserID:           userID,
		ServiceID:        serviceID,
		AggregatedAmount: valueobject.NewAmount(0),
	}
}

func (ca *ChargeAggregate) GetKey() string {
	return ca.UserID.String() + ":" + ca.ServiceID.String()
}
