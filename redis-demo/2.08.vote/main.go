package main

import (
	"fmt"

	"demo/redis-demo/2.08.vote/vote"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	v := vote.NewVote(client, "10086")

	up, err := v.Up("user1")
	if err != nil {
		return
	}
	fmt.Println(up)
	v.Up("user2")

	v.Down("user3")
	v.Down("user4")

	isVoted, err := v.IsVoted("user1")
	if err != nil {
		return
	}
	fmt.Println("Is voted user1", isVoted)

	count, err := v.UpCount()
	if err != nil {
		return
	}
	fmt.Println("Up count", count)

	downCount, err := v.DownCount()
	if err != nil {
		return
	}
	fmt.Println("Down count", downCount)

	total, err := v.Total()
	if err != nil {
		return
	}
	fmt.Println("Total count", total)
}
