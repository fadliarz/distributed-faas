package domain

type Invocation struct {
	InvocationID  InvocationID
	FunctionID    FunctionID
	UserID        UserID
	SourceCodeURL SourceCodeURL
	OutputURL     OutputURL
	Status        Status
	Timestamp     Timestamp
	IsRetry       IsRetry
}
