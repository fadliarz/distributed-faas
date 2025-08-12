package domain

import (
	"time"

	"github.com/fadliarz/distributed-faas/common/valueobject"
)

// Amount represents the charge amount in the smallest currency unit (e.g., cents)
type Amount int64

func NewAmount(amount int64) Amount {
	return Amount(amount)
}

func (a Amount) Int64() int64 {
	return int64(a)
}

// Charge represents a single charge event
type Charge struct {
	UserID    UserID
	ServiceID ServiceID
	Amount    Amount
	Timestamp valueobject.Timestamp
}

func NewCharge(userID UserID, serviceID ServiceID, amount Amount) *Charge {
	return &Charge{
		UserID:    userID,
		ServiceID: serviceID,
		Amount:    amount,
		Timestamp: valueobject.NewTimestamp(time.Now().Unix()),
	}
}

// ChargeAggregate represents aggregated charges for a specific user-service combination
type ChargeAggregate struct {
	UserID      UserID
	ServiceID   ServiceID
	TotalAmount Amount
	ChargeCount int64
	FirstCharge valueobject.Timestamp
	LastCharge  valueobject.Timestamp
}

func NewChargeAggregate(userID UserID, serviceID ServiceID) *ChargeAggregate {
	return &ChargeAggregate{
		UserID:      userID,
		ServiceID:   serviceID,
		TotalAmount: NewAmount(0),
		ChargeCount: 0,
		FirstCharge: valueobject.NewTimestamp(time.Now().Unix()),
		LastCharge:  valueobject.NewTimestamp(time.Now().Unix()),
	}
}

func (ca *ChargeAggregate) AddCharge(charge *Charge) {
	if ca.ChargeCount == 0 {
		ca.FirstCharge = charge.Timestamp
	}

	ca.TotalAmount = NewAmount(ca.TotalAmount.Int64() + charge.Amount.Int64())
	ca.ChargeCount++
	ca.LastCharge = charge.Timestamp
}

func (ca *ChargeAggregate) GetKey() string {
	return ca.UserID.String() + ":" + ca.ServiceID.String()
}
