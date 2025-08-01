package domain

type Checkpoint struct {
	CheckpointID  CheckpointID
	FunctionID    FunctionID
	UserID UserID
	SourceCodeURL SourceCodeURL
	Timestamp     Timestamp
	Status        Status
	OutputURL     OutputURL
	IsRetry       IsRetry
}
