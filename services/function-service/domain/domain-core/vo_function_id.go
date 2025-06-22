package domain

type FunctionID string

func (i *FunctionID) String() string {
	return string(*i)
}

func NewFunctionID(id string) (FunctionID, error) {
	return FunctionID(id), nil
}
