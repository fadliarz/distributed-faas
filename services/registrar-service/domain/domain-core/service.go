package domain

import "fmt"

type RegistrarDomainServiceImpl struct{}

func NewRegistrarDomainService() RegistrarDomainService {
	return &RegistrarDomainServiceImpl{}
}

func (s *RegistrarDomainServiceImpl) ValidateAndInitiateMachine(machine *Machine) error {
	// Validate the machine
	if machine.MachineID.String() != "" {
		return fmt.Errorf("machine ID is already initialized")
	}

	if machine.Address.String() == "" {
		return fmt.Errorf("machine address cannot be empty")
	}

	// Initiate the machine
	machine.Status = NewStatusFromInt(int(Available))

	return nil
}
