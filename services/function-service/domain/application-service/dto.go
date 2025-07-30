package application

type CreateFunctionCommand struct {
	UserID        string
	SourceCodeURL string
}

type GetUploadPresignedURLCommand struct {
	UserID     string
	FunctionID string
	Language   string
}
