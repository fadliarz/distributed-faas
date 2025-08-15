package application

type BillingCalculationEvent struct {
	UserID     string `json:"_id"`
	LastBilled int64  `json:"last_billed"`
}
