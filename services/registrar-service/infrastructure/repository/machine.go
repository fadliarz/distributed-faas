package repository

import (
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type MachineEntity struct {
	MachineID primitive.ObjectID `bson:"_id,omitempty"`
	Address   string             `bson:"address"`
	Status    string             `bson:"status"`
}

// Mapper

type MachineMapper interface {
	Entity(machine *domain.Machine) *MachineEntity
	Domain(entity *MachineEntity) *domain.Machine
}
