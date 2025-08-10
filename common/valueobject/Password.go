package valueobject

type Password string

func NewPassword(password string) Password {
	return Password(password)
}

func (p Password) String() string {
	return string(p)
}
