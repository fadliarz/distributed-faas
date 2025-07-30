package domain

type Address string

func NewAddress(url string) Address {
	return Address(url)
}

func (u *Address) String() string {
	return string(*u)
}
