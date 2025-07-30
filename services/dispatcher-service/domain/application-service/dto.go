package application

type InvocationCreatedEvent struct {
	InvocationID  string `json:"_id"`
	FunctionID    string `json:"function_id"`
	SourceCodeURL string `json:"source_code_url"`
	Timestamp     int64  `json:"timestamp"`
	IsRetry       bool   `json:"is_retry"`
}
