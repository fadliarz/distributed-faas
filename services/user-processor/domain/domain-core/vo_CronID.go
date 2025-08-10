package domain

type CronID string

func NewCronID(cronID string) CronID {
	return CronID(cronID)
}

func (c CronID) String() string {
	return string(c)
}
