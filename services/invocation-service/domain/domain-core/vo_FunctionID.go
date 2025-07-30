package domain

type FunctionID string

func NewFunctionID(id string) FunctionID {
	return FunctionID(id)
}

func (f FunctionID) String() string {
	return string(f)
}
