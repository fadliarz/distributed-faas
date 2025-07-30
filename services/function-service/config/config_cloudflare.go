package config

import (
	"fmt"
	"os"
)

type InputCloudflareConfig struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

func NewInputCloudflareConfig() (*InputCloudflareConfig, error) {
	accountID := os.Getenv("INPUT_CLOUDFLARE_ACCOUNT_ID")
	if accountID == "" {
		return nil, fmt.Errorf("missing INPUT_CLOUDFLARE_ACCOUNT_ID environment variable")
	}

	accessKeyID := os.Getenv("INPUT_CLOUDFLARE_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, fmt.Errorf("missing INPUT_CLOUDFLARE_ACCESS_KEY_ID environment variable")
	}

	secretAccessKey := os.Getenv("INPUT_CLOUDFLARE_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return nil, fmt.Errorf("missing INPUT_CLOUDFLARE_SECRET_ACCESS_KEY environment variable")
	}

	bucketName := os.Getenv("INPUT_CLOUDFLARE_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("missing INPUT_CLOUDFLARE_BUCKET_NAME environment variable")
	}

	return &InputCloudflareConfig{
		AccountID:       accountID,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
	}, nil
}
