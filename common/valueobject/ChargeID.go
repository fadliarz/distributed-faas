package valueobject

type ChargeID string

func NewChargeID(id string) ChargeID {
	return ChargeID(id)
}

func (c ChargeID) String() string {
	return string(c)
}
