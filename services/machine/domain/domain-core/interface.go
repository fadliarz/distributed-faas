package domain

type MachineDomainService interface {
	ValidateAndInitiateCheckpoint(checkpoint *Checkpoint, checkpointID CheckpointID) error
}
