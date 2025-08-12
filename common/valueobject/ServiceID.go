package valueobject

type ServiceID string

func NewServiceID(id string) ServiceID {
	return ServiceID(id)
}

func (s ServiceID) String() string {
	return string(s)
}
