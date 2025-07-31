package config

import (
	"fmt"
	"os"
	"strconv"
)

type RetryConfig struct {
	RetryIntervalInSec int64
	ThresholdInSec     int64
}

func NewRetryConfig() (*RetryConfig, error) {
	config := &RetryConfig{}

	var err error

	retryInterval := os.Getenv("RETRY_INTERVAL_IN_SEC")
	if retryInterval == "" {
		return nil, fmt.Errorf("RETRY_INTERVAL_IN_SEC environment variable is not set")
	}

	if config.RetryIntervalInSec, err = strconv.ParseInt(retryInterval, 10, 64); err != nil {
		return nil, fmt.Errorf("RETRY_INTERVAL_IN_SEC environment variable is not valid")
	}

	thresholdInSec := os.Getenv("THRESHOLD_IN_SEC")
	if thresholdInSec == "" {
		return nil, fmt.Errorf("THRESHOLD_IN_SEC environment variable is not set")
	}

	if config.ThresholdInSec, err = strconv.ParseInt(thresholdInSec, 10, 64); err != nil {
		return nil, fmt.Errorf("THRESHOLD_IN_SEC environment variable is not valid")
	}

	return config, nil
}
