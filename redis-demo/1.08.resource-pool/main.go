package main

import (
	"fmt"

	"demo/redis-demo/1.08.resource-pool/resource_pool"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	pool := resource_pool.NewRedisResourcePool(client, "free", "busy")

	_, err := pool.Associate("resource1")
	if err != nil {
		return
	}
	pool.Associate("resource2")
	pool.Associate("resource3")

	resource, _ := pool.Acquire()
	fmt.Println(resource)

	pool.Release("resource1")
	pool.Release("resource2")

	resource, _ = pool.Acquire()
	fmt.Println(resource)

	_, err2 := pool.Disassociate("resource3")
	if err2 != nil {
		return
	}

	has, err := pool.Has("resource3")
	if err != nil {
		return
	}
	fmt.Println(has)

	availCount, err := pool.AvailableCount()
	if err != nil {
		return
	}
	fmt.Println(availCount)

	occupiedCount, err := pool.OccupiedCount()
	if err != nil {
		return
	}
	fmt.Println(occupiedCount)
}
