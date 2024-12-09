package main

import (
	"fmt"

	"demo/redis-demo/2.06.session/session"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	ss := session.NewRedisSession(client, "10086")

	token, err := ss.Create()
	if err != nil {
		panic(err)
	}

	status, err := ss.Validate(token)
	if err != nil {
		panic(err)
	}
	if status != session.TokenValid {
		fmt.Println("token状态为", status, "校验失败")
	}

	ss.Destroy()

	status, err = ss.Validate(token)
	fmt.Println(status, err)
}
