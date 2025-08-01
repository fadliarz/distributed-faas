package domain

import (
	"fmt"
	"time"
)

type InvocationDomainServiceImpl struct{}

func NewInvocationDomainService() InvocationDomainService {
	return &InvocationDomainServiceImpl{}
}

func (s *InvocationDomainServiceImpl) ValidateAndInitiateInvocation(invocation *Invocation, invocationID string, function *Function) error {
	if invocation == nil {
		return fmt.Errorf("invocation cannot be nil")
	}

	if invocation.InvocationID != "" {
		return fmt.Errorf("invocation ID must be empty for a new invocation, got: %s", invocation.InvocationID)
	}

	if invocation.SourceCodeURL != "" {
		return fmt.Errorf("source code URL must be empty for a new invocation, got: %s", invocation.SourceCodeURL)
	}

	if invocation.Timestamp != 0 {
		return fmt.Errorf("timestamp must be zero for a new invocation, got: %d", invocation.Timestamp.Int64())
	}

	if invocation.IsRetry != false {
		return fmt.Errorf("isRetry must be false for a new invocation, got: %t", invocation.IsRetry.Bool())
	}

	if invocation.OutputURL != "" {
		return fmt.Errorf("output URL must be empty for a new invocation, got: %s", invocation.OutputURL.String())
	}

	if invocation.Status != 0 {
		return fmt.Errorf("status must be Unknown for a new invocation, got: %s", invocation.Status.String())
	}

	invocation.InvocationID = NewInvocationID(invocationID)
	invocation.SourceCodeURL = function.SourceCodeURL
	invocation.Status = NewStatus("PENDING")
	invocation.Timestamp = NewTimestamp(time.Now().Unix())
	invocation.IsRetry = NewIsRetry(false)

	if invocation.FunctionID == "" {
		return fmt.Errorf("function ID cannot be empty for invocation: %s", invocation.InvocationID)
	}

	return nil
}
