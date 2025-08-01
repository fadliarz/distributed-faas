package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fadliarz/distributed-faas/services/machine/config"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

// Constructors

type CommandHandler struct {
	service MachineApplicationService
}

func NewCommandHandler(service MachineApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
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

// Methods

func (h *CommandHandler) ProcessInvocation(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	// Persist checkpoint
	checkpointID, err := h.service.PersistCheckpoint(ctx, cmd)

	// Ignore if the checkpoint already exists
	if err != nil && errors.Is(err, domain.ErrCheckpointAlreadyExists) {
		return checkpointID, nil
	}

	// Ignore if the checkpoint has already been reprocessed
	if err != nil && errors.Is(err, domain.ErrCheckpointAlreadyReprocessed) {
		return checkpointID, nil
	}

	if err != nil {
		return "", fmt.Errorf("failed to persist checkpoint: %w", err)
	}

	go h.service.ExecuteFunction(context.TODO(), cmd.SourceCodeURL, cmd.FunctionID, cmd.InvocationID)

	return checkpointID, nil
}
