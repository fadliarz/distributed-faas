package domain

type Status int

const (
	Unknown Status = iota
	Pending
	Retrying
	Reprocessing
	Success
)

func (s Status) String() string {
	return []string{"UNKNOWN", "PENDING", "RETRYING", "REPROCESSING", "SUCCESS"}[s]
}

func NewStatus(status string) Status {
	hashMap := map[string]int{
		"PENDING":      1,
		"RETRYING":     2,
		"REPROCESSING": 3,
		"SUCCESS":      4,
	}

	return Status(hashMap[status])
}
