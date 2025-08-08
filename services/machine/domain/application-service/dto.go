package application

type ProcessInvocationCommand struct {
	InvocationID  string
	FunctionID    string
	UserID        string
	SourceCodeURL string
	Status        string
	Timestamp     int64
}
