package valueobject

type LastBilled int64

func NewLastBilled(lastBilled int64) LastBilled {
	return LastBilled(lastBilled)
}

func (lb *LastBilled) Int64() int64 {
	return int64(*lb)
}
