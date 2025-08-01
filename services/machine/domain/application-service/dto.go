package application

type ProcessInvocationCommand struct {
	InvocationID  string
	FunctionID    string
	UserID        string
	SourceCodeURL string
	Timestamp     int64
	IsRetry       bool
}
