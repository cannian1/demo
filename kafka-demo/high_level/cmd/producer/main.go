package main

import (
	"context"
	"demo/kafka-demo/high_level/internal"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/cast"
)

const (
	kafkaURL1  = "localhost:9092"
	kafkaURL2  = "localhost:9093"
	kafkaURL3  = "localhost:9094"
	kafkaTopic = "bare_recycle"
)

func main() {

	kafkaWriter := &kafka.Writer{
		Addr: kafka.TCP([]string{kafkaURL1, kafkaURL2, kafkaURL3}...),
		// Topic:                  kafkaTopic,          // 主题不能在 Writer 和 Message 里同时指定
		Balancer:               &kafka.LeastBytes{},    // 将消息发送到收到最少字节数的分区
		MaxAttempts:            5,                      // 重试次数，默认 10 次
		RequiredAcks:           kafka.RequireOne,       // 只要分区的 leader 副本成功写入消息，那么它就会收到来自服务端的成功响应。
		Async:                  false,                  // 是否异步发送，false 阻塞直到消息发送成功，如果发送失败，返回错误；如果设为 true，发送失败不会返回错误，在不关心是否发送成功时用。
		Completion:             nil,                    // 发送完成后的回调函数，如果 Async 为 true，那么 Completion 会在消息发送成功或失败后被调用。
		Compression:            kafka.Gzip,             // 消息压缩算法，默认不压缩
		Logger:                 nil,                    // 不为 nil 时，Writer 会使用它来记录日志
		ErrorLogger:            nil,                    // 不为 nil 时，Writer 会使用它来记录错误日志
		Transport:              kafka.DefaultTransport, // 用于发送消息的传输层，默认是 DefaultTransport，即连接超时 3s，空闲连接 30s
		AllowAutoTopicCreation: false,                  // 是否允许自动创建主题，如果为 true，当主题不存在时，Writer 会尝试创建主题，一般生产环境不建议开启
	}

	eventStreamWriter := internal.NewKafkaProducer(kafkaWriter)
	defer eventStreamWriter.Close()

	for i := range 3 {
		msg := producerHandler()

		err := eventStreamWriter.Produce(context.Background(), kafkaTopic, msg...)
		if err != nil {
			log.Printf("生产者%d出现错误:%v", i, err)
			return
		}
	}
	fmt.Println("完成")
}

func producerHandler() []internal.Message {
	uid, _ := uuid.NewV7()
	infos := []map[string]any{
		{
			"uid":             uid.String(),
			"zone":            3,
			"region":          1,
			"frozen_strategy": "7 Days",
		},
		{
			"uid":             uid.String(),
			"zone":            3,
			"region":          2,
			"frozen_strategy": "5 Days",
		},
		{
			"uid":             uid.String(),
			"zone":            2,
			"region":          2,
			"frozen_strategy": "14 Days",
		},
	}

	msgs := make([]internal.Message, 0, len(infos))
	for i, v := range infos {
		val, _ := json.Marshal(v)
		msgs = append(msgs, internal.Message{
			Key:   cast.ToString(i),
			Value: val,
		})
	}

	return msgs
}
