package unique_counter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type UniqueCounter interface {
	// Include 判断元素是否存在，不存在则增加计数器的值
	Include(string) (bool, error)
	// Exclude 判断元素是否存在，存在则减少计数器的值
	Exclude(string) (bool, error)
	// Count 获取计数器的值
	Count() (int64, error)
}

type RedisUniqueCounter struct {
	client redis.Cmdable
	key    string
}

func NewRedisUniqueCounter(client redis.Cmdable, key string) *RedisUniqueCounter {
	return &RedisUniqueCounter{client: client, key: key}
}

func (c *RedisUniqueCounter) Include(member string) (bool, error) {
	op, err := c.client.SAdd(context.Background(), c.key, member).Result()
	return op == 1, err
}

func (c *RedisUniqueCounter) Exclude(member string) (bool, error) {
	op, err := c.client.SRem(context.Background(), c.key, member).Result()
	return op == 1, err
}

func (c *RedisUniqueCounter) Count() (int64, error) {
	return c.client.SCard(context.Background(), c.key).Result()
}
