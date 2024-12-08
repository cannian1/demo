package lock

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 基于 setnx 和 del 实现的单点锁

// redLock 多主过半成功视为锁上，这里不实现

const valueOfLock = "*"

type RedisLock struct {
	client *redis.Client
	key    string
}

func NewRedisLock(client *redis.Client, key string) *RedisLock {
	return &RedisLock{client: client, key: key}
}

// Acquire 获取锁
func (l *RedisLock) Acquire() bool {
	return l.client.SetNX(context.Background(), l.key, valueOfLock, 0).Val()
}

// Release 释放锁
func (l *RedisLock) Release() bool {
	return l.client.Del(context.Background(), l.key).Val() == 1
}
