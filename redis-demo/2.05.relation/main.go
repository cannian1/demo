package main

import (
	"fmt"

	"demo/redis-demo/2.05.relation/relation"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	petter := relation.NewRedisRelation(client, "Peter")
	bob := relation.NewRedisRelation(client, "Bob")

	petter.Follow("Tom")
	petter.Follow("Jerry")

	bob.Follow("Tom")
	bob.Follow("Peter")

	pFollowers, _ := petter.FollowersCount()
	pFollowings, _ := petter.FollowingsCount()

	bFollowers, _ := bob.FollowersCount()
	bFollowings, _ := bob.FollowingsCount()

	fmt.Println("Peter's followers:", pFollowers)
	fmt.Println("Peter's followings:", pFollowings)
	fmt.Println("Bob's followers:", bFollowers)
	fmt.Println("Bob's followings:", bFollowings)

	isFriends, _ := petter.IsFriend("Bob")
	fmt.Println("Peter and Bob are friends:", isFriends)

	petter.Follow("Bob")
	isFriends, _ = petter.IsFriend("Bob")
	fmt.Println("Peter and Bob are friends:", isFriends)

	fmt.Println(petter.IsFollowing("Bob"))
}
