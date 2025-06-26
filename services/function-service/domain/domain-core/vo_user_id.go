package domain

type UserID string

func (i *UserID) String() string {
	return string(*i)
}

func NewUserID(id string) (UserID, error) {
	return UserID(id), nil
}
