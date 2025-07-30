package domain

type InvocationID string

func NewInvocationID(id string) InvocationID {
	return InvocationID(id)
}

func (i InvocationID) String() string {
	return string(i)
}
