package main

import (
	"fmt"

	"demo/redis-demo/4.counter/counter"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	var cnt counter.Counter
	cnt = counter.NewRedisCounter(client, "user:1:counter")

	err := cnt.Reset(200)
	if err != nil {
		fmt.Println(err)
		return
	}

	cnt.Incr()
	cnt.Incr()
	cnt.Incr()

	val, err := cnt.Get()
	fmt.Println(val, err)

	cnt.Decr()
	cnt.Decr()

	val, err = cnt.Get()
	fmt.Println(val, err)
	cnt.Reset()
	fmt.Println(cnt.Get())
}
