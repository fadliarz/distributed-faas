package domain

type InvocationDomainService interface {
	ValidateAndInitiateInvocation(invocation *Invocation, invocationID string, function *Function) error
}
