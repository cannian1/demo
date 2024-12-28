package internal

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(reader *kafka.Reader) *KafkaConsumer {
	return &KafkaConsumer{reader: reader}
}

func (c *KafkaConsumer) Consume(ctx context.Context) (<-chan Message, error) {
	messageCh := make(chan Message)

	go func() {
		defer close(messageCh)
		for {
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				break
			}
			messageCh <- Message{
				Key:   string(m.Key),
				Value: m.Value,
				Time:  m.Time,
			}
		}
	}()

	return messageCh, nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
