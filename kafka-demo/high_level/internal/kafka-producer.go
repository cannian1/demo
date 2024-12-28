package internal

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(writer *kafka.Writer) *KafkaProducer {
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) Produce(ctx context.Context, topic string, msgs ...Message) error {
	kafkaMsgs := make([]kafka.Message, 0, len(msgs))
	for _, v := range msgs {
		kafkaMsgs = append(kafkaMsgs, kafka.Message{
			Topic: topic,
			Key:   []byte(v.Key),
			Value: v.Value,
			//WriterData: nil,
		})
	}

	return p.writer.WriteMessages(ctx, kafkaMsgs...)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
