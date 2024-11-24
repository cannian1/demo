package main

import (
	"fmt"
	"time"
)

const (
	Small = iota
	Middle
	Big
	Supper
)

// 定义接口
type Prototype interface {
	Clone() Prototype
}

// FrenchFries 薯条 具体原型
type FrenchFries struct {
	Source        string
	Specification int
	PresetTime    time.Time
}

// 实现 Clone 方法
func (p *FrenchFries) Clone() Prototype {
	newPrototype := *p
	return &newPrototype
}

func main() {
	// 创建原型对象
	frenchFries := &FrenchFries{
		Source:        "tomato",
		Specification: Middle,
		PresetTime:    time.Date(2024, time.February, 21, 3, 4, 5, 0, time.UTC),
	}

	// 使用原型对象创建新对象
	newObject := frenchFries.Clone()

	myFrenchFries := newObject.(*FrenchFries)
	myFrenchFries.Source = "salt"
	myFrenchFries.Specification = Big

	fmt.Printf("%+v\n", myFrenchFries)
}
