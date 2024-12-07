package main

import (
	"encoding/json"
	"fmt"
	"time"

	"demo/redis-demo/1.cache/cache"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

const (
	id       = 10086            // 请求 id
	ttl      = 60 * time.Second // 最大超时时间
	reqTimes = 3                // 请求次数
)

type Person struct {
	ID     int64  `json:"id" redis:"id"` // hset 需要这个 tag
	Name   string `json:"name" redis:"name"`
	Gender string `json:"gender" redis:"gender"`
	Age    int64  `json:"age" redis:"age"`
}

func main() {
	var myCache cache.Cache[any]
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	myCache = cache.NewRedisCache(client)
	hashCache := cache.NewRedisHashCache(client)

	// 缓存示例 1
	for range reqTimes {
		post, _ := myCache.Get(cast.ToString(id))
		if len(post) == 0 {
			post = getPostFromTemplate(id)
			myCache.Set(cast.ToString(id), post, ttl)
			fmt.Println("缓存不存在，已从 db 和 template 获取数据并存入缓存")
		} else {
			fmt.Println("缓存存在，从缓存中可以获取到数据")
		}
	}

	p := &Person{
		ID:     1,
		Name:   "马Peter",
		Gender: "m",
		Age:    18,
	}

	jsonData, _ := json.Marshal(p)

	// 缓存示例 2 - JSON
	myCache.Set("person:10086", jsonData, ttl)
	jsonRes, _ := myCache.Get("person:10086")
	fmt.Println(jsonRes) // {"id":1,"name":"马Peter","gender":"m","age":18}

	var p2 Person
	json.Unmarshal([]byte(jsonRes), &p2)

	err := hashCache.Set("hash:person:10086", p, ttl)
	if err != nil {
		fmt.Println(err)
		return
	}

	mpRes, err := hashCache.Get("hash:person:10086")
	fmt.Println(mpRes, err)
}

// 模拟从数据库中获取数据
func getFromDB(id int) string {
	return "hello world"
}

// 模拟使用数据和模板生成HTML页面
func getPostFromTemplate(id int) string {
	content := getFromDB(id)
	return fmt.Sprintf("<html><p>%s</p></html>", content)
}
