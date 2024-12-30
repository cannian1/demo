# kafka-go 示例

## 安装依赖
```bash
go get github.com/segmentio/kafka-go
```

## 介绍
这个库有两套 API，这里只用高级别的API做示例。

## 注意事项
这个库目前不支持幂等生产者和事务的配置。

而且这个库的维护者不是很活跃，如果要用更多的特性，建议使用 [sarama](https://github.com/IBM/sarama) 或 [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go)。

- sarama: 更建议使用这个
- confluent-kafka-go: 这个库需要额外的依赖，并且要使用 CGO，不建议使用