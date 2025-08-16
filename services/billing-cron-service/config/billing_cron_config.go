package config

import (
	"fmt"
	"os"
	"strconv"
)

type BillingCronConfig struct {
	CronIntervalInSec int64
}

func NewBillingCronConfig() (*BillingCronConfig, error) {
	config := &BillingCronConfig{}

	var err error

	cronInterval := os.Getenv("CRON_INTERVAL_IN_SEC")
	if cronInterval == "" {
		return nil, fmt.Errorf("CRON_INTERVAL_IN_SEC environment variable is not set")
	}

	if config.CronIntervalInSec, err = strconv.ParseInt(cronInterval, 10, 64); err != nil {
		return nil, fmt.Errorf("CRON_INTERVAL_IN_SEC environment variable is not valid")
	}

	return config, nil
}
