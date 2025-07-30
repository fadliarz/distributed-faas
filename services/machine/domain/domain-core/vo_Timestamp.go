package domain

type Timestamp int64

func NewTimestamp(timestamp int64) Timestamp {
	return Timestamp(timestamp)
}

func (t Timestamp) Int64() int64 {
	return int64(t)
}
