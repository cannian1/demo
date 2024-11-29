package main

import "fmt"

// Mediator 中介者接口
type Mediator interface {
	sendMessage(msg string, user User)
	receiveMessage() string
}

// ChatRoom 具体中介者
type ChatRoom struct {
	Message string
}

func (c *ChatRoom) sendMessage(msg string, user User) {
	c.Message = fmt.Sprintf("%s 发送消息: %s\n", user.name, msg)
}

func (c *ChatRoom) receiveMessage() string {
	return c.Message
}

// User 用户类
type User struct {
	name     string
	mediator Mediator
}

func (u *User) getName() string {
	return u.name
}

func (u *User) setMediator(mediator Mediator) {
	u.mediator = mediator
}

func (u *User) sendMessage(msg string) {
	u.mediator.sendMessage(msg, *u)
}

func (u *User) receiveMessage() string {
	return u.mediator.receiveMessage()
}

func main() {
	// 创建中介者对象
	chatRoom := &ChatRoom{}

	// 创建用户对象，并设置中介者
	user1 := &User{name: "User1"}
	user2 := &User{name: "User2"}
	user1.setMediator(chatRoom)
	user2.setMediator(chatRoom)

	// 用户发送消息
	user1.sendMessage("Hello World!")
	fmt.Println(user2.receiveMessage())
	user2.sendMessage("Hi!")
	fmt.Println(user1.receiveMessage())
}
