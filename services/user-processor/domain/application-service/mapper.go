package application

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
)

// Constructors

type UserProcessorDataMapperImpl struct{}

func NewUserProcessorDataMapper() UserProcessorDataMapper {
	return &UserProcessorDataMapperImpl{}
}

func (m *UserProcessorDataMapperImpl) UserEventToCron(event *UserEvent) *domain.Cron {
	return &domain.Cron{
		UserID: valueobject.NewUserID(event.UserID),
	}
}
