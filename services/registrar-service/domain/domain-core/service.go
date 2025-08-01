package domain

import (
	"fmt"
)

type RegistrarDomainServiceImpl struct{}

func NewRegistrarDomainService() RegistrarDomainService {
	return &RegistrarDomainServiceImpl{}
}

func (s *RegistrarDomainServiceImpl) ValidateAndInitiateMachine(machine *Machine, machineID MachineID) error {
	if machine.MachineID.String() != "" {
		return fmt.Errorf("machine ID is already initialized")
	}

	if machine.Status != 0 {
		return fmt.Errorf("machine status is already set to %s", machine.Status.String())
	}

	// Initiate the machine
	machine.MachineID = NewMachineID(machineID.String())
	machine.Status = NewStatusFromInt(int(Available))

	if machine.Address.String() == "" {
		return fmt.Errorf("machine address cannot be empty")
	}

	return nil
}
