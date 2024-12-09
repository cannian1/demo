package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type HashCache[T any] interface {
	Set(key string, value T, ttl time.Duration) error
	Get(key string) (map[string]string, error)
}

type RedisHashCache struct {
	client *redis.Client
}

func NewRedisHashCache(client *redis.Client) *RedisHashCache {
	return &RedisHashCache{client: client}
}

func (cli *RedisHashCache) Set(key string, value any, ttl time.Duration) error {
	tx := cli.client.TxPipeline()
	ctx := context.Background()
	tx.HSet(ctx, key, value) // 如果是传入结构体的话，需要打上 redis 的 tag
	tx.Expire(ctx, key, ttl)

	if _, err := tx.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (cli *RedisHashCache) Get(key string) (map[string]string, error) {
	return cli.client.HGetAll(context.Background(), key).Result()
}
