package main

import (
	"fmt"

	"demo/redis-demo/2.09.rank/rank"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	v := rank.NewRank(client, "peak")

	v.SetWeight("小米YU7", 100)
	v.SetWeight("组建叙利亚过度政府", 90)
	v.SetWeight("南方游客刚到漠河", 80)

	length, err := v.Length()
	if err != nil {
		return
	}
	fmt.Println("排行榜长度", length)

	weight, err := v.GetWeight("小米YU7")
	if err != nil {
		return
	}
	fmt.Println("小米YU7的权重", weight)

	top, err := v.Top(2)
	if err != nil {
		return
	}
	fmt.Println("排行榜前两名", top)

	bottom, err := v.Bottom(2)
	if err != nil {
		return
	}
	fmt.Println("排行榜后两名", bottom)

	v.UpdateWeight("小米YU7", -20)

	weight, err = v.GetWeight("小米YU7")
	if err != nil {
		return
	}
	fmt.Println("小米YU7的权重", weight)
}
