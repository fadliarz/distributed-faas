package repository

type InvocationEntity struct {
	FunctionID   string `bson:"function_id"`
	InvocationID string `bson:"invocation_id"`
}
