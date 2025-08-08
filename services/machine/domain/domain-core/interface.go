package domain

type MachineDomainService interface {
	ValidateAndInitiateCheckpoint(checkpoint *Checkpoint) error
}
