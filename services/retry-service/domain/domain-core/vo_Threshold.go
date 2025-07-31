package domain

type Threshold int64

func NewThreshold(threshold int64) Threshold {
	return Threshold(threshold)
}

func (t Threshold) Int64() int64 {
	return int64(t)
}
