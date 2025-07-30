package domain

type InvocationDomainServiceImpl struct{}

func NewInvocationDomainService() InvocationDomainService {
	return &InvocationDomainServiceImpl{}
}
