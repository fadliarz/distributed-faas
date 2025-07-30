package domain

import "fmt"

type MachineDomainServiceImpl struct{}

func NewMachineDomainService() MachineDomainService {
	return &MachineDomainServiceImpl{}
}

func (s *MachineDomainServiceImpl) ValidateAndInitiateCheckpoint(checkpoint *Checkpoint, checkpointID CheckpointID) error {
	if checkpoint.OutputURL.String() != "" {
		return fmt.Errorf("output URL must be empty for a new checkpoint")
	}

	if checkpoint.Status != 0 {
		return fmt.Errorf("checkpoint status must be zero (uninitialized)")
	}
	checkpoint.Status = NewStatusFromInt(int(Pending))

	if checkpoint.FunctionID.String() == "" {
		return fmt.Errorf("function ID cannot be empty")
	}

	if checkpoint.SourceCodeURL.String() == "" {
		return fmt.Errorf("source code URL cannot be empty")
	}

	return nil
}
