package valueobject

type UserID string

func NewUserID(id string) UserID {
	return UserID(id)
}

func (u UserID) String() string {
	return string(u)
}
