package main

import (
	"fmt"
	"log"

	"demo/redis-demo/1.10.iterator/dbiterator"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 创建 DbIterator 实例
	iterator := dbiterator.NewDbIterator(client, 5)

	// 使用迭代器遍历数据库
	for {
		keys, err := iterator.Next()
		if err != nil {
			log.Fatalf("Error during iteration: %v", err)
		}

		if len(keys) == 0 {
			// 迭代完成
			break
		}

		fmt.Println("Keys:", keys)
	}

	// 重置迭代器
	iterator.Rewind()

	// 再次迭代
	for {
		keys, err := iterator.Next()
		if err != nil {
			log.Fatalf("Error during iteration: %v", err)
		}

		if len(keys) == 0 {
			break
		}

		fmt.Println("Keys after rewind:", keys)
	}
}
