package main

import (
	"demo/copier-demo/copier"
	"fmt"
	"time"
)

type A struct {
	Name  string
	Age   int
	CTime int64
	Time  time.Time
}

type B struct {
	Name  string
	Age   int
	Ctime time.Time
	Time  int64
}

func main() {
	a := A{
		Name:  "Name-A",
		Age:   11,
		CTime: time.Now().Unix(),
		Time:  time.Date(2025, 3, 21, 18, 0, 0, 0, time.UTC),
	}
	var b B

	_ = copier.CopyWithOption(&b, &a, copier.Int64ToTime(), copier.TimeToInt64())
	fmt.Printf("%v", b)
}
