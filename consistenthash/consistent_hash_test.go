package consistenthash

import (
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"sort"
	"testing"
)

// Benchmark_SelectNode-12    	10000000	       117 ns/op	       0 B/op	       0 allocs/op  2*12
// Benchmark_SelectNode-12    	10000000	       154 ns/op	       0 B/op	       0 allocs/op  12*12
// Benchmark_SelectNode-12    	10000000	       204 ns/op	       0 B/op	       0 allocs/op  30*12
// Benchmark_SelectNode-12    	 5000000	       258 ns/op	       0 B/op	       0 allocs/op  50*12
// Benchmark_SelectNode-12    	 3000000	       404 ns/op	       0 B/op	       0 allocs/op  100*12
// Benchmark_SelectNode-12    	 1000000	      1195 ns/op	       0 B/op	       0 allocs/op  200*12
// Benchmark_SelectNode-12    	  500000	      2667 ns/op	       0 B/op	       0 allocs/op  500*12
func Benchmark_SelectNode(b *testing.B) {
	h := New(30)
	h.AddNodes(ipFactory(2)...)
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		h.SelectNode("172.16.1.1")
	}
}

// 2*12
// Min 3rd probability: 0.094% 0.100% 0.101%
// Max 3rd probability: 52.029% 19.926% 10.760%
// Avg probability: 4.167%

// 12*12
// Min 3rd probability: 0.042% 0.042% 0.043%
// Max 3rd probability: 5.387% 5.066% 4.460%
// Avg probability: 0.694%

// 30*12
// Min 3rd probability: 0.004% 0.008% 0.008%
// Max 3rd probability: 1.323% 1.299% 1.141%
// Avg probability: 0.278%

// 50*12
// Min 3rd probability: 0.002% 0.004% 0.004%
// Max 3rd probability: 0.783% 0.706% 0.625%
// Avg probability: 0.167%

// 100*12
// Min 3rd probability: 0.000% 0.000% 0.001%
// Max 3rd probability: 0.353% 0.352% 0.346%
// Avg probability: 0.083%
func TestSelectNode(t *testing.T) {
	h := New(500)
	h.AddNodes(ipFactory(12)...)
	h.Show()
	h.SelectNode("172.16.1.1")
}

func TestHashLoop_AddNodes(t *testing.T) {
	type args struct {
		ips []string
	}
	tests := []struct {
		name string
		h    *HashLoop
		args args
	}{
		{name: "normal0", h: New(3), args: args{[]string{"1", "2", "3", "4"}}},
		{name: "normal1", h: New(6), args: args{[]string{"1", "2", "3", "4"}}},
		{name: "normal2", h: New(9), args: args{[]string{"1", "2", "3", "4"}}},
		{name: "normal3", h: New(12), args: args{[]string{"1", "2", "3", "4"}}},
		{name: "normal4", h: New(15), args: args{[]string{"1", "2", "3", "4"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.AddNodes(tt.args.ips...)
			if len(tt.h.nodes) != tt.h.virtualNum*len(tt.args.ips) {
				t.Fail()
			}
		})
	}
}

func newh() *HashLoop {
	h := New(3)
	h.AddNodes("1", "2", "3", "4")
	return h
}

func TestHashLoop_Select(t *testing.T) {
	h := New(200)
	h.AddNodes(
		"192.168.0.2",
		"192.168.0.3",
		"192.168.0.4",
		"192.168.0.5",
		"192.168.0.6",
		"192.168.0.7",
		"192.168.0.8",
		"192.168.0.9",
	)
	ips := ipFactory(1000000)
	res := make(map[string]int)
	for _, ip := range ips {
		r := h.SelectNode(ip)
		if _, ok := res[r]; ok {
			res[r]++
		} else {
			res[r] = 1
		}
	}
	for ip, count := range res {
		fmt.Println(ip, count)
	}
}

func TestHashLoop_DelNodes(t *testing.T) {
	type args struct {
		ips []string
	}
	tests := []struct {
		name string
		h    *HashLoop
		args args
	}{
		{name: "normal0", h: newh(), args: args{[]string{"1", "2", "3", "4"}}},
		{name: "normal0", h: newh(), args: args{[]string{"1", "2", "3"}}},
		{name: "normal0", h: newh(), args: args{[]string{"1", "2"}}},
		{name: "normal0", h: newh(), args: args{[]string{"1"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.DelNodes(tt.args.ips...)
			if len(tt.h.nodes) != 12-(tt.h.virtualNum*len(tt.args.ips)) {
				t.Fatal(fmt.Sprintf("want %d, get %d", len(tt.h.nodes), 12-(tt.h.virtualNum*len(tt.args.ips))))
			}
			tt.h.Show()
		})
	}
}

// Benchmark_crc32-12    	50000000	        22.2 ns/op	       0 B/op	       0 allocs/op
func Benchmark_crc32(b *testing.B) {
	c := crc32.NewIEEE()
	data := []byte("192.168.1.1")
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		c.Write(data)
		c.Sum32()
		c.Reset()
	}
}

// Benchmark_fnv-12    	200000000	         7.67 ns/op	       0 B/op	       0 allocs/op
func Benchmark_fnv(b *testing.B) {
	c := fnv.New32a()
	data := []byte("192.168.1.1")
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		c.Write(data)
		c.Sum32()
		c.Reset()
	}
}

func TestSort(t *testing.T) {
	var ns = nodeSet{
		node{ip: "1", num: 10},
		node{ip: "1", num: 20},
	}
	fmt.Println(ns)
	ns = append(ns, node{ip: "1", num: 15})
	fmt.Println(ns)
	sort.Sort(ns)
	fmt.Println(ns)
}
