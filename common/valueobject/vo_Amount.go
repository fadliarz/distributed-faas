package valueobject

type Amount int64

func NewAmount(amount int64) Amount {
	return Amount(amount)
}

func (a Amount) Int64() int64 {
	return int64(a)
}

func (a Amount) Add(other Amount) Amount {
	return Amount(int64(a) + int64(other))
}
