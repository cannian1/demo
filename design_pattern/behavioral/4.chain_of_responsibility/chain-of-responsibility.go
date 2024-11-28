package main

import "fmt"

// Handler 是职责链中的处理者接口
type Handler interface {
	Handle(request string) bool
	SetNext(handler Handler) // 设置下一个处理者
}

// ConcreteHandler 是一个具体的处理者
type ConcreteHandler struct {
	next Handler
	name string
}

func (c *ConcreteHandler) Handle(request string) bool {
	// 判断当前处理者是否能处理请求
	if c.canHandle(request) {
		fmt.Printf("%s is handling request: %s\n", c.name, request)
		return true
	}

	// 如果当前处理者不能处理，传递给下一个处理者
	if c.next != nil {
		return c.next.Handle(request)
	}

	// 如果没有处理者能够处理请求
	fmt.Println("Request could not be handled")
	return false
}

func (c *ConcreteHandler) SetNext(handler Handler) {
	c.next = handler
}

// canHandle 是具体处理者处理请求的逻辑
func (c *ConcreteHandler) canHandle(request string) bool {
	// 在这里实现处理的条件，简单演示为处理特定请求
	return request == c.name
}

func main() {
	// 创建处理者
	handlerA := &ConcreteHandler{name: "Handler A"}
	handlerB := &ConcreteHandler{name: "Handler B"}
	handlerC := &ConcreteHandler{name: "Handler C"}

	// 设置职责链
	handlerA.SetNext(handlerB)
	handlerB.SetNext(handlerC)

	// 测试请求
	requests := []string{"Handler A", "Handler B", "Handler C", "Handler D"}
	for _, request := range requests {
		fmt.Printf("Processing request: %s\n", request)
		handlerA.Handle(request)
		fmt.Println()
	}
}
