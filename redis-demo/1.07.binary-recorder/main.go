package main

import (
	"demo/redis-demo/1.07.binary-recorder/binary_record"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	binaryRecord := binary_record.NewBinaryRecord(client, "user:1:sign_in")

	binaryRecord.SetBit(1)
	binaryRecord.SetBit(3)

	count, _ := binaryRecord.CountBit(0, 6)
	println(count)
}
