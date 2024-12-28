package main

import (
	"context"
	"demo/kafka-demo/high_level/internal"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaURL1  = "localhost:9092"
	kafkaURL2  = "localhost:9093"
	kafkaURL3  = "localhost:9094"
	kafkaTopic = "bare_recycle"
)

func main() {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaURL1, kafkaURL2, kafkaURL3},
		GroupID:     "bare-consumerGroup-1", // 消费者组 ID，如果这个字段被指定，Partition 分区就不应该被指定了
		GroupTopics: nil,                    // 指定多个主题用
		Topic:       kafkaTopic,
		// Partition:              0, // 分区，和 GroupID 互斥
		Dialer:           kafka.DefaultDialer, // 默认拨号器，超时时间 10s，拨号器还可以配置 KeepAlive、事务ID等
		QueueCapacity:    100,                 // 队列容量，默认 100
		MinBytes:         1,                   // 默认1。最少从 broker 收到一批多少个字节，如果低容量主题设置这个参数很大，可能导致延迟交付
		MaxBytes:         1e6,                 // 默认 1MB。最多多少字节，超过则被截断
		MaxWait:          10 * time.Second,    // 默认 10s。如果没有足够的数据可用，最多等待的时间。如果不满一批，等待这个时间后返回当前可用数据
		ReadBatchTimeout: 10 * time.Second,    // 默认 10s。Kafka 的消息批中提取单个消息的超时时间
		ReadLagInterval:  1 * time.Second,     // 更新消费者“滞后”（Lag）信息的频率，设为负数则取消上报。
		// 消费者组的分区分配策略。当多个消费者属于同一个消费组时，Kafka 会根据分区分配策略决定如何将主题的分区分配给这些消费者
		GroupBalancers: []kafka.GroupBalancer{ // 里面的顺序即策略的优先顺序
			kafka.RangeGroupBalancer{},      // 按照范围分配
			kafka.RoundRobinGroupBalancer{}, // 按照轮询分配 (最均衡)
			// kafka.RackAffinityGroupBalancer{}, // 机架亲和分配 (最低时延)
		},
		HeartbeatInterval:      3 * time.Second,  // 默认3s。心跳
		CommitInterval:         0,                // 向代理提交偏移量的时间间隔。如果为0，提交将被同步处理.(仅当 GroupID 被设置时生效)
		PartitionWatchInterval: 5 * time.Second,  // 默认5s。reader 多久检查分区被改变。如果被改变了，例如添加了新的分区，则会进行 rebalance。(GroupID 和 WatchPartitionChanges 都被设置时生效)
		WatchPartitionChanges:  true,             // 用于通知 kafka-go，如果主题发生任何分区更改，消费者组应该轮询协调器并重新平衡
		SessionTimeout:         30 * time.Second, // 会话超时时间。在协调器认为消费者已死并启动再平衡之前没有心跳的时间长度
		RebalanceTimeout:       30 * time.Second, // 默认30s。设置协调器等待成员加入的时间长度，作为再平衡的一部分。对于负载较高的 kafka 服务器，把这个值设置的比较高更好(仅当 GroupID 被设置时生效)。
		JoinGroupBackoff:       5 * time.Second,  // 默认5s。出现错误后重新加入消费者组的等待时间。
		RetentionTime:          12 * time.Hour,   // 默认24h。消费者组将被 broker 保留的时间(仅当 GroupID 被设置时生效)

		// 当发现一个没有提交偏移量的分区时，StartOffset决定消费者组应该从哪里开始消费。如果非零，则必须设置为 FirstOffset 或 LastOffset 中的一个。
		StartOffset:           kafka.FirstOffset,     // (仅当 GroupID 被设置时生效)
		ReadBackoffMin:        0,                     // 默认 100ms。reader 轮询新消息的最小等待时间
		ReadBackoffMax:        0,                     // 默认 1s。reader 轮询新消息的最小等待时间
		Logger:                nil,                   // 日志
		ErrorLogger:           nil,                   // 错误日志
		IsolationLevel:        kafka.ReadUncommitted, // 隔离级别。读未提交级别下所有记录都可见。读提交级别下仅非事务的记录和已提交的记录可见
		MaxAttempts:           3,                     // 默认3次。返回错误前的最大重试次数
		OffsetOutOfRangeError: false,                 // 向后兼容的字段，未来版本会移除。目的是出现 OffsetOutOfRange 错误时直接返回错误而不是再重试
	})

	eventStreamReader := internal.NewKafkaConsumer(kafkaReader)
	defer eventStreamReader.Close()

	messages, err := eventStreamReader.Consume(context.Background())
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range messages {
		fmt.Printf("【Received message】key:%v val:%v,time:%v\n", msg.Key, string(msg.Value), msg.Time)
	}
}
