package domain

type UserProcessorDomainService interface {
	ValidateAndInitiateCron(cron *Cron) error
}
