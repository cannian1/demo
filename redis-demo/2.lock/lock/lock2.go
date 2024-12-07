package lock

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const maxExpireTime = 5 * time.Minute

type RedisLock2 struct {
	client *redis.Client
	key    string
}

func NewRedisLock2(client *redis.Client, key string) *RedisLock2 {
	return &RedisLock2{client: client, key: key}
}

// Acquire 获取锁
func (l *RedisLock2) Acquire(owner string, timeout time.Duration) bool {
	if timeout > maxExpireTime {
		timeout = maxExpireTime
	}
	return l.client.SetNX(context.Background(), l.key, owner, timeout).Val()
}

// Renew 续期锁
func (l *RedisLock2) Renew(owner string, timeout time.Duration) bool {
	if timeout > maxExpireTime {
		timeout = maxExpireTime
	}

	ctx := context.Background()

	err := l.client.Watch(ctx, func(tx *redis.Tx) error {

		currentOwner, err := tx.Get(ctx, l.key).Result()
		if err != nil {
			log.Println("续期阶段 Error checking lock owner:", err)
			return err
		}

		// 如果当前持有者不是自己，不能续期
		if currentOwner != owner {
			log.Println("你不是该锁当前的持有者，无权续期")
			return fmt.Errorf("not the lock owner\n")
		}

		return tx.Expire(ctx, l.key, timeout).Err()
	})

	if err != nil {
		log.Println("Error renewing lock:", err)
		return false
	}
	return true
}

// Release 释放锁
func (l *RedisLock2) Release(owner string) bool {
	ctx := context.Background()

	err := l.client.Watch(ctx, func(tx *redis.Tx) error {

		currentOwner, err := tx.Get(ctx, l.key).Result()
		if err != nil {
			log.Println("释放阶段 Error checking lock owner:", err)
			return err
		}

		// 如果当前持有者不是自己，不能续期
		fmt.Println("currentOwner:", currentOwner)
		fmt.Println("owner:", owner)
		if currentOwner != owner {
			log.Println("你不是该锁当前的持有者，无权释放")
			return fmt.Errorf("not the lock owner\n")
		}

		return tx.Del(ctx, l.key).Err()
	})

	if err != nil {
		log.Println("Error release lock:", err)
		return false
	}
	return true
}
