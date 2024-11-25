package main

import "fmt"

// Flyweight 接口定义享元对象的方法
type Flyweight interface {
	Operation(extrinsicState string)
}

// ConcreteFlyweight 是具体的享元类，包含内部状态和外部状态
type ConcreteFlyweight struct {
	intrinsicState string // 内部状态
}

// Operation 实现 Flyweight 接口中定义的操作方法
func (f *ConcreteFlyweight) Operation(extrinsicState string) {
	fmt.Printf("ConcreteFlyweight: intrinsic state is %s, extrinsic state is %s\n", f.intrinsicState, extrinsicState)
}

// FlyweightFactory 是享元工厂类，负责创建和管理共享对象
type FlyweightFactory struct {
	flyweights map[string]Flyweight
}

// GetFlyweight 从享元工厂中获取共享对象
func (factory *FlyweightFactory) GetFlyweight(key string) Flyweight {
	if flyweight, ok := factory.flyweights[key]; ok {
		return flyweight
	}
	flyweight := &ConcreteFlyweight{intrinsicState: key}
	factory.flyweights[key] = flyweight
	return flyweight
}

func main() {
	factory := &FlyweightFactory{flyweights: make(map[string]Flyweight)}

	// 获取共享对象，如果工厂中不存在则创建
	flyweight1 := factory.GetFlyweight("flyweight1")
	flyweight2 := factory.GetFlyweight("flyweight2")
	flyweight3 := factory.GetFlyweight("flyweight3")
	flyweight4 := factory.GetFlyweight("flyweight4")
	flyweight5 := factory.GetFlyweight("flyweight5")

	// 调用共享对象的操作方法
	flyweight1.Operation("external state 1")
	flyweight2.Operation("external state 2")
	flyweight3.Operation("external state 3")
	flyweight4.Operation("external state 4")
	flyweight5.Operation("external state 5")
}
