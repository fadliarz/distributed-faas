package application

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
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

func (s *MachineApplicationServiceImpl) ExecuteFunction(ctx context.Context, url string, functionID, invocationID string) error {
	log.Debug().Msgf("Executing function with ID %s and invocation ID %s", functionID, invocationID)

	sourceCode, err := s.fetchCodeFromURL(ctx, url)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to fetch code from URL")

		return fmt.Errorf("failed to fetch code from URL: %w", err)
	}

	logs, err := s.executeCode(ctx, sourceCode)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to execute code")

		return fmt.Errorf("failed to execute code: %w", err)
	}

	err = s.uploadLog(ctx, logs, functionID, invocationID)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to upload logs")

		return fmt.Errorf("failed to upload logs: %w", err)
	}

	log.Debug().Msg("Function execution completed successfully")

	return nil
}

func (s *MachineApplicationServiceImpl) fetchCodeFromURL(ctx context.Context, url string) (string, error) {
	resp, err := s.client.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.config.Cloudflare.BucketName,
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

	return string(contentBytes), nil
}

func (s *MachineApplicationServiceImpl) executeCode(ctx context.Context, code string) ([]LogLine, error) {
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

		if err := scanner.Err(); err != nil {
			log.Warn().Err(err).Msg("Error reading stderr")
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	cmd.Wait()
	wg.Wait()

	if len(logLines) < 1 {
		return nil, nil
	}

	return logLines, nil
}

func (s *MachineApplicationServiceImpl) uploadLog(ctx context.Context, logs []LogLine, functionID, invocationID string) error {
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

	_, err = s.client.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.config.Cloudflare.BucketName,
		Key:    &key,
		Body:   bytes.NewReader(jsonBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to upload execution result to S3: %w", err)
	}

	err = s.repositoryManager.Checkpoint.UpdateStatusToSuccess(ctx, domain.NewCheckpointID(invocationID), domain.NewOutputURL(key))
	if err != nil {
		// Send critical log if the update fails

		log.Warn().Err(err).Msgf("Failed to update checkpoint status to SUCCESS for ID %s", invocationID)
	}

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
