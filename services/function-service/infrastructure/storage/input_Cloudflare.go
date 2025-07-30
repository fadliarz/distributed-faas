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
	client *s3.Client

	config *config.InputCloudflareConfig
}

func (s *InputCloudflareStorage) GetUploadPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	request, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
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
