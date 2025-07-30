package domain

import "fmt"

type Address string

func NewAddress(address string) (Address, error) {
	if address == "" {
		return "", fmt.Errorf("address cannot be empty")
	}

	return Address(address), nil
}

func NewLooseAddress(address string) Address {
	return Address(address)
}

func (u *Address) String() string {
	return string(*u)
}
