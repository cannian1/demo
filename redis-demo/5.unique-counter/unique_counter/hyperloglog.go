package unique_counter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 用于不需要非常精确的计数场景，比如统计UV，HyperLogLog 可以实现基数统计
// 优点：占用空间小，时间复杂度低

type RedisUniqueCounterHyperLogLog struct {
	client *redis.Client
	key    string
}

func NewRedisUniqueCounterHyperLogLog(client *redis.Client, key string) *RedisUniqueCounterHyperLogLog {
	return &RedisUniqueCounterHyperLogLog{client: client, key: key}
}

func (c *RedisUniqueCounterHyperLogLog) Include(member string) (bool, error) {
	op, err := c.client.PFAdd(context.Background(), c.key, member).Result()
	return op == 1, err
}

func (c *RedisUniqueCounterHyperLogLog) Exclude(member string) (bool, error) {
	return false, nil
}

func (c *RedisUniqueCounterHyperLogLog) Count() (int64, error) {
	return c.client.PFCount(context.Background(), c.key).Result()
}
