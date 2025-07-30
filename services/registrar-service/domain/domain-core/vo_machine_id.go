package domain

import "fmt"

type MachineID string

func NewMachineID(id string) (MachineID, error) {
	if id == "" {
		return "", fmt.Errorf("machine ID cannot be empty")
	}

	return MachineID(id), nil
}

func NewLooseMachineID(id string) MachineID {
	return MachineID(id)
}

func (f MachineID) String() string {
	return string(f)
}
