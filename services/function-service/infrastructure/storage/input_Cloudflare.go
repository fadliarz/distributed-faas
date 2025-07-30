package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fadliarz/distributed-faas/services/function-service/config"
)

type InputCloudflareStorage struct {
	config *config.InputCloudflareConfig

	client    *s3.Client
	presigner *s3.PresignClient
}

func NewInputCloudflareStorage(config *config.InputCloudflareConfig, client *s3.Client) *InputCloudflareStorage {
	return &InputCloudflareStorage{
		config:    config,
		client:    client,
		presigner: s3.NewPresignClient(client),
	}
}

func (s *InputCloudflareStorage) GetUploadPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	request, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("failed to create presigned upload URL: %w", err)
	}

	return request.URL, nil
}
