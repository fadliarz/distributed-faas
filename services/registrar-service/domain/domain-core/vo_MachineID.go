package domain

type MachineID string

func NewMachineID(id string) MachineID {

	return MachineID(id)
}

func (f MachineID) String() string {
	return string(f)
}
