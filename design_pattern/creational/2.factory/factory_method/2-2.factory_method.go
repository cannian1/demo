package main

import "fmt"

var ActIDNameMap = map[int64]string{
	1: "慕君之意望君觉",
	2: "如梦而逝的崇高者们",
}

// ActivityConfig 活动配置
type ActivityConfig interface {
	Show()
}

type Factory interface {
	CreateActivityConfig() ActivityConfig
}

type DoubleNinja struct {
	ID   int64
	Name string
}

func (dn DoubleNinja) Show() {
	fmt.Printf("活动1 %s\n", dn.Name)
}

func (dn DoubleNinja) CreateActivityConfig() ActivityConfig {
	var id int64 = 1
	return DoubleNinja{
		ID:   id,
		Name: ActIDNameMap[id],
	}
}

type WhiteTiger struct {
	ActivityConfig
	ID   int64
	Name string
}

func (wt WhiteTiger) Show() {
	fmt.Printf("活动2 %s\n", wt.Name)
}

func (wt WhiteTiger) CreateActivityConfig() ActivityConfig {
	var id int64 = 2
	return &WhiteTiger{
		ID:   id,
		Name: ActIDNameMap[id],
	}
}

func main() {
	var doubleNinjaFac Factory
	doubleNinjaFac = DoubleNinja{} // 赋值为具体的实现 (适合依赖注入)
	d := doubleNinjaFac.CreateActivityConfig()
	d.Show()

	var whiteTigerFac Factory
	whiteTigerFac = WhiteTiger{} // 赋值为具体的实现 (适合依赖注入)
	w := whiteTigerFac.CreateActivityConfig()
	w.Show()
}
