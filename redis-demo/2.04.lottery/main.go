package main

import (
	"fmt"

	"demo/redis-demo/2.04.lottery/lottery"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	l := lottery.NewRedisLottery(client, "my_lottery")

	l.Join("player1")
	l.Join("player2")
	l.Join("player3")

	size, _ := l.Size()
	fmt.Println("Size:", size)

	players, _ := l.Draw(2, false)
	fmt.Println("Players:", players)

	size, _ = l.Size()
	fmt.Println("Size:", size)

	players, _ = l.Draw(2, true)
	fmt.Println("Players:", players)

	size, _ = l.Size()
	fmt.Println("Size:", size)
}
