package domain

type FunctionDomainService interface {
	ValidateAndInitiateFunction(function *Function, functionID string) error
}
