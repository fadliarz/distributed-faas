package application

type CreateChargeCommand struct {
	UserID    string
	ServiceID string
	Amount    int64
}
