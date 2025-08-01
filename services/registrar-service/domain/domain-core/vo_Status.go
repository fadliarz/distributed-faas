package domain

type Status int

const (
	Unknown Status = iota
	Available
	Unavailable
)

func (s Status) String() string {
	return []string{"Unknown", "Available", "Unavailable"}[s]
}

func NewStatus(status string) Status {
	hashMap := map[string]int{
		"Available":   1,
		"Unavailable": 2,
	}

	return Status(hashMap[status])
}

func NewStatusFromInt(status int) Status {
	return Status(status)
}
