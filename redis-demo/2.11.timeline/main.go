package main

import (
	"fmt"

	"demo/redis-demo/2.11.timeline/timeline"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	tl := timeline.NewPaging(client, "TopicTimeline")

	tl.Add([]timeline.Content{
		{
			Item:   "topic:10086",
			Weight: 1748707200,
		},
		{
			Item:   "topic:10087",
			Weight: 1748704120,
		},
		{
			Item:   "topic:10001",
			Weight: 1748701586,
		},
		{
			Item:   "topic:10084",
			Weight: 1748733067,
		},
		{
			Item:   "topic:10072",
			Weight: 1748712345,
		},
	}...)

	get, err := tl.Get(1, 5)
	if err != nil {
		return
	}
	fmt.Println(get)

	getWithTime, err := tl.GetWithTime(1, 5)
	if err != nil {
		return
	}
	fmt.Println(getWithTime)

	getByTimeRange, err := tl.GetByTimeRange(1748700000, 1748720000, 1, 5)
	if err != nil {
		return
	}
	fmt.Println(getByTimeRange)

	length, err := tl.Length()
	if err != nil {
		return
	}
	fmt.Println("Length", length)

	number, err := tl.Number(5)
	if err != nil {
		return
	}
	fmt.Println("Number", number)
}
