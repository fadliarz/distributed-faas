package domain

type UserID string

func NewUserID(id string) UserID {
	return UserID(id)
}

func (f UserID) String() string {
	return string(f)
}
