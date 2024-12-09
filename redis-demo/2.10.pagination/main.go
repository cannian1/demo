package main

import (
	"fmt"

	"demo/redis-demo/2.10.pagination/paging"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	p := paging.NewPaging(client, "TopicList")

	for i := range 10 {
		added, _ := p.Add("Topic:" + cast.ToString(i))
		fmt.Println("Added", added)
	}

	length, _ := p.Length()
	fmt.Println("Length", length)

	page1, _ := p.Get(1, 3)
	page2, _ := p.Get(2, 3)
	page3, _ := p.Get(3, 3)

	fmt.Println("Page 1", page1)
	fmt.Println("Page 2", page2)
	fmt.Println("Page 3", page3)
}
