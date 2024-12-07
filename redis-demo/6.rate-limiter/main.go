package main

import (
	"time"

	"demo/redis-demo/6.rate-limiter/ratelimit"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	rateLimiter := ratelimit.NewRedisRateLimiter(client, "publish", 60, 10)

	for range 11 {
		time.Sleep(time.Millisecond * 500)
		allowed, err := rateLimiter.IsAllowed(1)
		println(allowed, err)
	}

	duration, err := rateLimiter.Duration(1)
	println(duration, err)

	remaining, err := rateLimiter.Remaining(1)
	println(remaining, err)

	rateLimiter.Revoke(1)
}
