package mq

import "context"

// MQ is interface of message queue
type MQ interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Consume(ctx context.Context, topic string) ([]byte, error)
}
