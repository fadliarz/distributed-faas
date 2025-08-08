package domain

type Checkpoint struct {
	CheckpointID  CheckpointID
	FunctionID    FunctionID
	UserID        UserID
	SourceCodeURL SourceCodeURL
	Status        Status
	Timestamp     Timestamp
	OutputURL     OutputURL
}
