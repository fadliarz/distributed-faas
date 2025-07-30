package kafka

import "context"

type MessageProcessor[M any] interface {
	Process(ctx context.Context, msg M) error
}

type MessageDeserializer[M any] interface {
	Deserialize(data []byte) (M, error)
}

type ProcessException struct {
	Err         error
	IsRetryable bool
}
