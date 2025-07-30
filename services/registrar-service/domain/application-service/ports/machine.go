package ports

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

type MachineRepository interface {
	Save(ctx context.Context, machine *domain.Machine) (domain.MachineID, error)
	UpdateStatus(ctx context.Context, machineID domain.MachineID, address domain.Address, status domain.Status) error
}
