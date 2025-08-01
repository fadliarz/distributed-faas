package domain

type Status int

const (
	Unknown Status = iota
	Pending
	Completed
)

func (s Status) String() string {
	return []string{"UNKNOWN", "PENDING", "SUCCESS"}[s]
}

func NewStatus(status string) Status {
	hashMap := map[string]int{
		"PENDING":   1,
		"SUCCESS": 2,
	}

	return Status(hashMap[status])
}

func NewStatusFromInt(status int) Status {
	return Status(status)
}
