package main

import (
	"fmt"
)

var ActIDNameMap = map[int64]string{
	1: "慕君之意望君觉",
	2: "如梦而逝的崇高者们",
}

// ActivityConfig 活动配置
type ActivityConfig interface {
	Show()
	GetName() string
}

type DoubleNinja struct {
	ActivityConfig
	ID   int64
	Name string
}

func (dn DoubleNinja) GetName() string {
	return dn.Name
}

func (dn DoubleNinja) Show() {
	fmt.Printf("活动1 %s\n", dn.GetName())
}

type WhiteTiger struct {
	ActivityConfig
	ID   int64
	Name string
}

func (dn WhiteTiger) GetName() string {
	return dn.Name
}

func (dn WhiteTiger) Show() {
	fmt.Printf("活动2 %s\n", dn.GetName())
}

type Factory struct {
}

func (f Factory) CreateActivityConfig(id int64) ActivityConfig {
	switch id {
	case 1:
		return &DoubleNinja{
			ID:   id,
			Name: ActIDNameMap[id],
		}
	case 2:
		return &WhiteTiger{
			ID:   id,
			Name: ActIDNameMap[id],
		}
	default:
		return nil
	}
}

func main() {
	factory := Factory{}

	act1 := factory.CreateActivityConfig(1)
	act1.Show()

	act2 := factory.CreateActivityConfig(2)
	act2.Show()
}
