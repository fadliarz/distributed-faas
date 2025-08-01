package domain

type RegistrarDomainService interface {
	ValidateAndInitiateMachine(machine *Machine, machineID MachineID) error
}
