package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

type InputStorageImpl struct {
	cloduflare *InputCloudflareStorage
}

func NewInputStorage(cloudflare *InputCloudflareStorage) application.InputStorage {
	return &InputStorageImpl{
		cloduflare: cloudflare,
	}
}

func (s *InputStorageImpl) GetFunctionUploadPresignedURL(ctx context.Context, userID domain.UserID, functionID domain.FunctionID, language domain.Language, expiration time.Duration) (string, error) {
	url, err := s.cloduflare.GetUploadPresignedURL(ctx, s.getFunctionUploadKey(userID, functionID, language), expiration)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL for function upload: %w", err)
	}

	return url, nil
}

func (s *InputStorageImpl) getFunctionUploadKey(userID domain.UserID, functionID domain.FunctionID, language domain.Language) string {
	return fmt.Sprintf("%s/%s/main_%d.%s", userID.String(), functionID.String(), time.Now().Unix(), language.String())
}
