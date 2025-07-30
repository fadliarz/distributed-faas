package application

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"github.com/rs/zerolog/log"
)

type LogLine struct {
	Timestamp time.Time
	Content   string
	Source    string
}

func (h *CommandHandler) ProcessInvocation(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	// Persist checkpoint
	checkpointID, err := h.service.PersistCheckpoint(ctx, cmd)

	log.Debug().Err(err).Msg("Inspecting error after persisting checkpoint")

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

	log.Debug().Msgf("Executing function in separate goroutine")
	go h.executeFunction(context.TODO(), cmd.SourceCodeURL, cmd.FunctionID, cmd.InvocationID)

	return checkpointID, nil
}

func (h *CommandHandler) executeFunction(ctx context.Context, url string, functionID, invocationID string) error {
	// Fetch the code from the URL
	log.Debug().Msgf("Fetching code from URL: %s", url)
	sourceCode, err := h.fetchCodeFromURL(ctx, url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch code from URL")

		return fmt.Errorf("failed to fetch code from URL: %w", err)
	}

	// Execute the code
	log.Debug().Msg("Executing code")
	logs, err := h.executeCode(ctx, sourceCode)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute code")

		return fmt.Errorf("failed to execute code: %w", err)
	}

	log.Debug().Msgf("Execution completed with %d log lines", len(logs))
	// Upload the logs to S3
	err = h.uploadLog(ctx, logs, functionID, invocationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload logs")

		return fmt.Errorf("failed to upload logs: %w", err)
	}

	return nil
}

func (h *CommandHandler) fetchCodeFromURL(ctx context.Context, url string) (string, error) {
	resp, err := h.client.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &h.config.Cloudflare.BucketName,
		Key:    &url,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch code from URL: %w", err)
	}

	defer resp.Body.Close()
	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debug().Msgf("Fetched code content: %s", string(contentBytes))

	return string(contentBytes), nil
}

func (h *CommandHandler) executeCode(ctx context.Context, code string) ([]LogLine, error) {
	// Prepare the command to run the code in a Docker container
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"run", "--rm", "--network", "none", "-i",
		"--memory", "128m", "--cpus", "0.5",
		"node:18-alpine", "node", "-e", code)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var logLines []LogLine

	// Start goroutines to read stdout and stderr
	log.Debug().Msg("Starting goroutine to read stdout and stderr")

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			logLines = append(logLines, LogLine{
				Timestamp: time.Now(),
				Content:   line,
				Source:    "stdout",
			})
			mu.Unlock()
		}

		log.Debug().Msg("Finished reading stdout")

		if err := scanner.Err(); err != nil {
			log.Warn().Err(err).Msg("Error reading stdout")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			logLines = append(logLines, LogLine{
				Timestamp: time.Now(),
				Content:   line,
				Source:    "stderr",
			})

			mu.Unlock()
		}

		log.Debug().Msg("Finished reading stderr")

		if err := scanner.Err(); err != nil {
			log.Warn().Err(err).Msg("Error reading stderr")
		}
	}()

	log.Debug().Msgf("Executing command")

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	cmd.Wait()
	wg.Wait()

	log.Debug().Msg("Command execution completed")

	if len(logLines) < 1 {
		return nil, nil
	}

	return logLines, nil
}

func (h *CommandHandler) uploadLog(ctx context.Context, logs []LogLine, functionID, invocationID string) error {
	// Aggregate logs into a single string
	aggregatedLogs := interleaveLogs(logs)

	payload := map[string]interface{}{
		"logs": aggregatedLogs,
	}

	// Convert the payload to JSON
	jsonBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// Upload the JSON to S3
	key := fmt.Sprintf("%s/%s.json", functionID, invocationID)

	log.Debug().Msgf("Uploading execution result to S3 with key: %s", key)

	_, err = h.client.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &h.config.Cloudflare.BucketName,
		Key:    &key,
		Body:   bytes.NewReader(jsonBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to upload execution result to S3: %w", err)
	}

	log.Debug().Msgf("Execution result uploaded successfully to S3 with key: %s", key)

	return nil
}

func interleaveLogs(logLines []LogLine) string {
	// Sort log lines by timestamp
	sort.Slice(logLines, func(i, j int) bool {
		return logLines[i].Timestamp.Before(logLines[j].Timestamp)
	})

	// Combine log lines into a single string
	var combinedLogBuilder strings.Builder
	for _, ll := range logLines {
		combinedLogBuilder.WriteString(fmt.Sprintf("[%s] %s\n",
			ll.Timestamp.Format(time.RFC3339Nano), ll.Content))
	}

	return combinedLogBuilder.String()
}
