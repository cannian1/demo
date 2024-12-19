package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache[T any] interface {
	Set(key string, value T, ttl time.Duration) error
	Get(key string) (string, error)
}

type RedisCache struct {
	client redis.Cmdable
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{client: client}
}

func (cli *RedisCache) Set(key string, value any, ttl time.Duration) error {
	return cli.client.Set(context.Background(), key, value, ttl).Err()
}

func (cli *RedisCache) Get(key string) (string, error) {
	return cli.client.Get(context.Background(), key).Result()
}
