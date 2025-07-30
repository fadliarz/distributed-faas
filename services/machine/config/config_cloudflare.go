package config

import (
	"fmt"
	"os"
)

type OutputCloudflareConfig struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

func NewOutputCloudflareConfig() (*OutputCloudflareConfig, error) {
	accountID := os.Getenv("OUTPUT_CLOUDFLARE_ACCOUNT_ID")
	if accountID == "" {
		return nil, fmt.Errorf("missing OUTPUT_CLOUDFLARE_ACCOUNT_ID environment variable")
	}

	accessKeyID := os.Getenv("OUTPUT_CLOUDFLARE_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, fmt.Errorf("missing OUTPUT_CLOUDFLARE_ACCESS_KEY_ID environment variable")
	}

	secretAccessKey := os.Getenv("OUTPUT_CLOUDFLARE_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return nil, fmt.Errorf("missing OUTPUT_CLOUDFLARE_SECRET_ACCESS_KEY environment variable")
	}

	bucketName := os.Getenv("OUTPUT_CLOUDFLARE_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("missing OUTPUT_CLOUDFLARE_BUCKET_NAME environment variable")
	}

	return &OutputCloudflareConfig{
		AccountID:       accountID,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
	}, nil
}
