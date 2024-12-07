package main

import (
	"fmt"
	"sync"
	"time"

	"demo/redis-demo/2.lock/lock"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	redisLock := lock.NewRedisLock(client, "my-lock")
	redisLockWithExpire := lock.NewRedisLock2(client, "my-lock2")

	common(redisLock)
	withExpire(redisLockWithExpire)
}

func common(redisLock *lock.RedisLock) {
	fmt.Println("-----基础的基于 redis 实现的锁----")
	var wg sync.WaitGroup

	for i := range 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if !redisLock.Acquire() {
				fmt.Printf("%d 获取锁失败\n", i)
			} else {
				fmt.Printf("%d 加锁成功\n", i)
			}
		}(i)
	}

	wg.Wait()

	if !redisLock.Release() {
		fmt.Println("释放锁失败")
	} else {
		fmt.Println("释放锁成功")
	}
}

func withExpire(redisLockWithExpire *lock.RedisLock2) {
	fmt.Println("-----带超时、续期和身份认证的基于 redis 实现的锁----")

	owner1 := uuid.New().String()
	owner2 := uuid.New().String()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if !redisLockWithExpire.Acquire(owner1, time.Second) {
			fmt.Printf("%s 获取锁失败\n", owner1)
		} else {
			fmt.Printf("%s 加锁成功\n", owner1)
		}

		if !redisLockWithExpire.Renew(owner1, time.Second) {
			fmt.Printf("%s 续期失败\n", owner1)
		} else {
			fmt.Printf("%s 续期成功\n", owner1)
		}
	}()

	time.Sleep(1 * time.Second)

	// owner2 必然获取不到锁 ，因为 owner1 续期成功
	go func() {
		defer wg.Done()
		if !redisLockWithExpire.Acquire(owner2, 3*time.Second) {
			fmt.Printf("%s 获取锁失败\n", owner2)
		} else {
			fmt.Printf("%s 加锁成功\n", owner2)
		}
	}()

	wg.Wait()

	// 此时 owner1 已经超时释放锁，owner2 也没有获取到锁
	if !redisLockWithExpire.Release(owner1) {
		fmt.Printf("%s 释放锁失败\n", owner1)
	} else {
		fmt.Printf("%s 释放锁成功\n", owner1)
	}
}
