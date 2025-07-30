package domain

type FunctionID string

func (i *FunctionID) String() string {
	return string(*i)
}

func NewFunctionID(id string) FunctionID {
	return FunctionID(id)
}
