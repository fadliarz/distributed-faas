package domain

type Checkpoint struct {
	CheckpointID  CheckpointID
	FunctionID    FunctionID
	SourceCodeURL SourceCodeURL
	Timestamp     Timestamp
	Status        Status
	OutputURL     OutputURL
	IsRetry       IsRetry
}
