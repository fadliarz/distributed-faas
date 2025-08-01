package domain

type Address string

func NewAddress(address string) Address {
	return Address(address)
}

func (u *Address) String() string {
	return string(*u)
}
