package application

type CreateMachineCommand struct {
	Address string
}

type UpdateMachineStatusCommand struct {
	MachineID string
	Address   string
}
