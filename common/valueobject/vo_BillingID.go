package valueobject

type BillingID string

func NewBillingID(id string) BillingID {
	return BillingID(id)
}

func (b BillingID) String() string {
	return string(b)
}
