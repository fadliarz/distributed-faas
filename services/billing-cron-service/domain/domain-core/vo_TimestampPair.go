package domain

import "time"

// TimestampPair represents the before and after timestamps for LastBilled updates
type TimestampPair struct {
	Before int64
	After  int64
}

func NewTimestampPair() TimestampPair {
	now := time.Now()

	return TimestampPair{
		Before: time.Date(now.Year(), now.Month()-2, 0, 0, 0, 0, 0, time.UTC).Unix(),
		After:  time.Date(now.Year(), now.Month()-1, 0, 0, 0, 0, 0, time.UTC).Unix(),
	}
}

func (tp TimestampPair) BeforeTimestamp() int64 {
	return tp.Before
}

func (tp TimestampPair) AfterTimestamp() int64 {
	return tp.After
}
