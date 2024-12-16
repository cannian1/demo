package main

import (
	"fmt"

	"demo/redis-demo/2.01.bloom-filter/bloom_filter"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	// 创建布隆过滤器
	bf := bloom_filter.NewBloomFilterService(client, "my_bloom_filter", 1000, 0.1)
	err := bf.Add("hello")
	if err != nil {
		fmt.Println(err)
		return
	}
	bf.Add("world")

	// 检查是否存在
	exists, _ := bf.Exists("hello")
	if exists {
		fmt.Println("hello exists")
	} else {
		fmt.Println("hello not exists")
	}

	exists, _ = bf.Exists("calculate")
	if exists {
		fmt.Println("calculate exists")
	} else {
		fmt.Println("calculate not exists")
	}
}
