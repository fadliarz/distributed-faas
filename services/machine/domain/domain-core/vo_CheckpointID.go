package domain

type CheckpointID string

func NewCheckpointID(id string) CheckpointID {
	return CheckpointID(id)
}

func (c CheckpointID) String() string {
	return string(c)
}
