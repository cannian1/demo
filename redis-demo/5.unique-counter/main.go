package main

import (
	"demo/redis-demo/5.unique-counter/unique_counter"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	userCounter := unique_counter.NewRedisUniqueCounter(client, "visitCounter")
	userCounterHyperLogLog := unique_counter.NewRedisUniqueCounterHyperLogLog(client, "visitCounter2")
	userCounter.Include("user1")
	userCounter.Include("user2")
	userCounter.Include("user3")
	userCounter.Include("user2")

	count, _ := userCounter.Count()

	println(count)

	userCounterHyperLogLog.Include("user1")
	userCounterHyperLogLog.Include("user2")
	userCounterHyperLogLog.Include("user3")
	userCounterHyperLogLog.Include("user2")

	count, _ = userCounterHyperLogLog.Count()
	println(count)
}
