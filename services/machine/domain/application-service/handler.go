package application

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fadliarz/distributed-faas/services/machine/config"
)

// Constructors

type CommandHandler struct {
	service MachineApplicationService

	config *CommandHandlerConfig
	client *CommandHandlerClient
}

type CommandHandlerConfig struct {
	Cloudflare config.OutputCloudflareConfig
}

type CommandHandlerClient struct {
	S3 *s3.Client
}

func NewCommandHandler(service MachineApplicationService, config *CommandHandlerConfig, client *CommandHandlerClient) *CommandHandler {
	return &CommandHandler{
		service: service,
		config:  config,
		client:  client,
	}
}

func NewCommandHandlerConfig(cfg config.OutputCloudflareConfig) *CommandHandlerConfig {
	return &CommandHandlerConfig{
		Cloudflare: cfg,
	}
}

func NewCommandHandlerClient(s3 *s3.Client) *CommandHandlerClient {
	return &CommandHandlerClient{
		S3: s3,
	}
}
