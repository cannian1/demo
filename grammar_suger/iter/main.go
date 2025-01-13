package main

import (
	"fmt"
	"iter"
)

// Set 基于 map 定义一个存放元素的集合类型
type Set[E comparable] struct {
	m map[E]struct{}
}

// NewSet 返回一个 Set
func NewSet[E comparable]() *Set[E] {
	return &Set[E]{m: make(map[E]struct{})}
}

// Add 向 Set 中添加元素
func (s *Set[E]) Add(v E) {
	s.m[v] = struct{}{}
}

// All 返回一个迭代器，迭代集合中的所有元素
func (s *Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

func forRangeSet() {
	s := NewSet[string]()
	s.Add("Golang")
	s.Add("Java")
	s.Add("Python")
	s.Add("C++")

	for v := range s.All() {
		fmt.Println(v)
	}
}

func main() {
	forRangeSet()
}
