package domain

type Status int

const (
	Unknown Status = iota
	Pending
	Completed
)

func (s Status) String() string {
	return []string{"Unknown", "Pending", "Completed"}[s]
}

func NewStatus(status string) Status {
	hashMap := map[string]int{
		"Pending":   1,
		"Completed": 2,
	}

	return Status(hashMap[status])
}

func NewStatusFromInt(status int) Status {
	return Status(status)
}
