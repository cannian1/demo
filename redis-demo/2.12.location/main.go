package main

import (
	"fmt"

	"demo/redis-demo/2.12.location/location"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	tl := location.NewRedisLocation(client, "UserLocation")

	tl.Pin("Peter", 116.404145671154, 39.915874512598)
	tl.Pin("Tom", 116.404145678154, 39.915874512038)
	tl.Pin("Jerry", 116.404145678174, 39.915874512198)
	tl.Pin("John", 116.404145678156, 39.915874517090)
	tl.Pin("Alice", 116.40414562155, 39.915874512192)

	longitude, latitude, err := tl.Locate("Peter")
	if err != nil {
		return
	}
	fmt.Println("Peter's location is", longitude, latitude)

	distance, err := tl.Distance("Peter", "Tom")
	if err != nil {
		return
	}
	fmt.Println("The distance between Peter and Tom is", distance)

	size, err := tl.Size()
	if err != nil {
		return
	}
	fmt.Println("The number of users' location is", size)

	users, err := tl.Search("Peter", 10)
	if err != nil {
		return
	}
	fmt.Println("The users around Peter are", users)
}
