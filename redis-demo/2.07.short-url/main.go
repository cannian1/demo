package main

import (
	"fmt"

	"demo/redis-demo/2.07.short-url/shorturl"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	su := shorturl.NewRedisShortURL(client)

	su2 := shorturl.NewRedisShortURLWithCache(client)

	shortID, err := su.Shorten("https://github.com")
	if err != nil {
		panic(err)
	}

	fmt.Println("Shortened URL:", shortID)

	url, err := su.Restore(shortID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Restored URL:", url)

	cachedShortID, err := su2.Shorten("https://github.com")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Shortened URL with cache:", cachedShortID)

	url, err = su2.Restore(cachedShortID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Restored URL with cache:", url)
}
