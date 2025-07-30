package application

type CreateInvocationCommand struct {
	UserID     string
	FunctionID string
}

type GetInvocationQuery struct {
	UserID       string
	InvocationID string
}
