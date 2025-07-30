package storage

import (
	"context"
	"time"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service"
)

type InputStorageImpl struct {
	cloduflare *InputCloudflareStorage
}

func NewInputStorage(cloudflare *InputCloudflareStorage) application.InputStorage {
	return &InputStorageImpl{
		cloduflare: cloudflare,
	}
}

func (s *InputStorageImpl) GetUploadPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return s.cloduflare.GetUploadPresignedURL(ctx, key, expiration)
}
