package id_generator

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type IDGenerator interface {
	// Produce 生成 ID
	Produce() (int64, error)
	// Reserve 保留指定数量的前置 ID
	Reserve(int64) error
}

type RedisIDGenerator struct {
	client redis.Cmdable
	name   string
}

func NewRedisIDGenerator(client redis.Cmdable, name string) *RedisIDGenerator {
	return &RedisIDGenerator{client: client, name: name}
}

func (g *RedisIDGenerator) Produce() (int64, error) {
	// 还可以使用 HIncrBy 方法，把一组相关的 ID 生成器放在同一个 key 下管理
	return g.client.Incr(context.Background(), g.name).Result()
}

func (g *RedisIDGenerator) Reserve(id int64) error {
	return g.client.SetNX(context.Background(), g.name, id, 0).Err()
}
