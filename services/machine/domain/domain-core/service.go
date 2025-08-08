package domain

import "fmt"

type MachineDomainServiceImpl struct{}

func NewMachineDomainService() MachineDomainService {
	return &MachineDomainServiceImpl{}
}

func (s *MachineDomainServiceImpl) ValidateAndInitiateCheckpoint(checkpoint *Checkpoint) error {
	if checkpoint == nil {
		return fmt.Errorf("checkpoint cannot be nil")
	}

	if checkpoint.Status != Pending {
		return fmt.Errorf("checkpoint status must be 'Pending' for a new checkpoint")
	}

	if checkpoint.OutputURL.String() != "" {
		return fmt.Errorf("output URL must be empty for a new checkpoint")
	}

	checkpoint.Status = Pending

	if checkpoint.CheckpointID.String() == "" {
		return fmt.Errorf("checkpoint ID cannot be empty")
	}

	if checkpoint.FunctionID.String() == "" {
		return fmt.Errorf("function ID cannot be empty")
	}

	if checkpoint.UserID.String() == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	if checkpoint.SourceCodeURL.String() == "" {
		return fmt.Errorf("source code URL cannot be empty")
	}

	if checkpoint.Timestamp.Int64() <= 0 {
		return fmt.Errorf("checkpoint timestamp must be a positive integer")
	}

	return nil
}
