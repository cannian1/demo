package counter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Counter 计数器
type Counter interface {
	// Incr 增加计数器的值
	Incr(...int64) error
	// Decr 减少计数器的值
	Decr(...int64) error
	// Get 获取计数器的值
	Get() (int64, error)
	// Reset 重置计数器的值
	Reset(...int64) error
}

// 也可以用 HIncrBy 方法，把一组相关的计数器放在同一个 key 下管理
// HIncrBy 负数用于减少计数器的值

type RedisCounter struct {
	client *redis.Client
	key    string
}

func NewRedisCounter(client *redis.Client, key string) *RedisCounter {
	return &RedisCounter{client: client, key: key}
}

func (c *RedisCounter) Incr(incr ...int64) error {
	var delta int64 = 1
	if len(incr) > 0 {
		delta = incr[0]
	}
	return c.client.IncrBy(context.Background(), c.key, delta).Err()
}

func (c *RedisCounter) Decr(decr ...int64) error {
	var delta int64 = 1
	if len(decr) > 0 {
		delta = decr[0]
	}
	return c.client.DecrBy(context.Background(), c.key, delta).Err()
}

func (c *RedisCounter) Get() (int64, error) {
	return c.client.Get(context.Background(), c.key).Int64()
}

func (c *RedisCounter) Reset(reset ...int64) error {
	var value int64 = 0
	if len(reset) > 0 {
		value = reset[0]
	}
	return c.client.Set(context.Background(), c.key, value, 0).Err()
}
