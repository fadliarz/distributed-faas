package application

type CreateFunctionCommand struct {
	UserID string
}

type GetFunctionUploadPresignedURLQuery struct {
	UserID     string
	FunctionID string
	Language   string
}

type UpdateFunctionSourceCodeURLCommand struct {
	UserID        string
	FunctionID    string
	SourceCodeURL string
}
