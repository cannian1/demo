package main

import (
	"fmt"

	"demo/redis-demo/3.incr/id_generator"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	var gen id_generator.IDGenerator
	gen = id_generator.NewRedisIDGenerator(client, "my-id-gen")

	err := gen.Reserve(100000)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := gen.Produce()
	fmt.Println(id, err)

	id, err = gen.Produce()
	fmt.Println(id, err)

	id, err = gen.Produce()
	fmt.Println(id, err)
}
