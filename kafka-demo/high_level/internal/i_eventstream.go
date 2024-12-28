package internal

import (
	"context"
	"time"
)

type Message struct {
	Key   string
	Value []byte
	Time  time.Time
}

// EventStreamProducer 事件流生产者
type EventStreamProducer interface {
	Produce(ctx context.Context, topic string, msg ...Message) error
	Close() error
}

// EventStreamConsumer 事件流消费者
type EventStreamConsumer interface {
	Consume(ctx context.Context) (<-chan Message, error)
	Close() error
}
