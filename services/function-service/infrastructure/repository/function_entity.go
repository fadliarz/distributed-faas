package repository

type FunctionEntity struct {
	UserID        string `bson:"user_id"`
	FunctionID    string `bson:"function_id"`
	SourceCodeURL string `bson:"source_code_url"`
}
