package application

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fadliarz/distributed-faas/services/machine/config"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

// Constructors

type MachineApplicationServiceImpl struct {
	mapper            MachineDataMapper
	service           domain.MachineDomainService
	repositoryManager *MachineApplicationServiceRepositoryManager

	config *CommandHandlerConfig
	client *CommandHandlerClient
}

type MachineApplicationServiceRepositoryManager struct {
	Checkpoint CheckpointRepository
}

type CommandHandlerConfig struct {
	Cloudflare config.OutputCloudflareConfig
}

type CommandHandlerClient struct {
	S3 *s3.Client
}

func NewMachineApplicationService(
	mapper MachineDataMapper, service domain.MachineDomainService, repositoryManager *MachineApplicationServiceRepositoryManager,
	config *CommandHandlerConfig, client *CommandHandlerClient,
) MachineApplicationService {
	return &MachineApplicationServiceImpl{
		mapper:            mapper,
		service:           service,
		repositoryManager: repositoryManager,
		config:            config,
		client:            client,
	}
}

func NewMachineApplicationServiceRepositoryManager(checkpoint CheckpointRepository) *MachineApplicationServiceRepositoryManager {
	return &MachineApplicationServiceRepositoryManager{
		Checkpoint: checkpoint,
	}
}
