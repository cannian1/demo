package main

import (
	"fmt"
	"log"

	"demo/redis-demo/1.09.compact-string/compact_str"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	// 创建 CompactString 实例
	cs := compact_str.NewCompactStr(client, "my_compact_string", "\n")

	// 添加字符串
	l, err := cs.Append("HTTP/1.1 200 OK")
	if err != nil {
		log.Fatalf("Append failed: %v", err)
	}
	fmt.Println("Length:", l)

	l, err = cs.Append("Server: nginx/1.16.1")
	if err != nil {
		log.Fatalf("Append failed: %v", err)
	}
	fmt.Println("Length:", l)

	// 获取所有字符串
	strings, err := cs.GetBytes(compact_str.OptBytesRange{
		Start: 0,
		End:   -1,
	})
	if err != nil {
		log.Fatalf("GetBytes failed: %v", err)
	}

	// 输出结果
	fmt.Println("Strings:", strings)
}
