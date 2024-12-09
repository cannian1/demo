package main

import (
	"fmt"
	"time"

	"demo/redis-demo/2.03.auto-complete/auto_complete"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	ac := auto_complete.NewRedisAutoComplete(client, "my_auto_complete")

	ac.FeedEx("赛博朋克", 1, 10*time.Second)

	ac.Set("赛博朋克2077", 10)
	ac.Set("赛博朋克边缘行者", 9)
	ac.Set("赛博朋克2077捏脸", 8)
	ac.Set("赛博朋克酒保行动", 7)
	ac.Set("赛博朋克2077往日之影", 6)
	ac.Set("赛博朋克2077mod", 5)
	ac.Set("赛博朋克2077bgm", 4)
	ac.Set("赛博朋克捏脸", 3)
	ac.Set("赛博朋克2077不朽武器", 2)
	ac.Set("赛博朋克2077画面设置", 1)

	segments, err := ac.Hint("赛博朋", 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("Segments:", segments)
}
