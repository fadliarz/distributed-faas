package domain

type RegistrarDomainService interface {
	ValidateAndInitiateMachine(machine *Machine) error
}
