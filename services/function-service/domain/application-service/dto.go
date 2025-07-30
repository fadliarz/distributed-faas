package application

type CreateFunctionCommand struct {
	UserID        string
	SourceCodeURL string
}

type GetFunctionUploadPresignedURLQuery struct {
	UserID     string
	FunctionID string
	Language   string
}
