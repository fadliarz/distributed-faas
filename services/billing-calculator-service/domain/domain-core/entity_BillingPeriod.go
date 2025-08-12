package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type BillingPeriod struct {
	StartTimestamp int64
	EndTimestamp   int64
	UserID         valueobject.UserID
}
