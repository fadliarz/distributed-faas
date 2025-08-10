package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
)

// Constructors

type UserEventHandler struct {
	mapper            UserProcessorDataMapper
	domainService     domain.UserProcessorDomainService
	repositoryManager *UserEventHandlerRepositoryManager
}

type UserEventHandlerRepositoryManager struct {
	Cron CronRepository
}

func NewUserEventHandler(mapper UserProcessorDataMapper, domainService domain.UserProcessorDomainService, repositoryManager *UserEventHandlerRepositoryManager) *UserEventHandler {
	return &UserEventHandler{
		mapper:            mapper,
		domainService:     domainService,
		repositoryManager: repositoryManager,
	}
}

func NewUserEventHandlerRepositoryManager(cron CronRepository) *UserEventHandlerRepositoryManager {
	return &UserEventHandlerRepositoryManager{
		Cron: cron,
	}
}

// Methods

func (eh *UserEventHandler) HandleUserEvent(ctx context.Context, event *UserEvent) error {
	cron := eh.mapper.UserEventToCron(event)

	err := eh.domainService.ValidateAndInitiateCron(cron)
	if err != nil {
		return fmt.Errorf("failed to validate and initiate cron: %w", err)
	}

	_, err = eh.repositoryManager.Cron.Save(ctx, cron)
	if err != nil {
		return fmt.Errorf("failed to save cron: %w", err)
	}

	return nil
}
