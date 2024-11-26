package main

import "fmt"

// 抽象类，制作饮料,包裹一个模板的全部实现步骤
type Beverage interface {
	DoOnBeforeBeverage(fn func() error)
	BoilWater() //煮开水
	Brew()      //冲泡
	PourInCup() //倒入杯中
	AddThings() //添加酌料

	WantAddThings() bool //是否加入酌料Hook
}

// 封装一套流程模板，让具体的制作流程继承且实现
type template struct {
	b               Beverage
	onBeverageStart func() error
}

// 封装的固定模板
func (t *template) MakeBeverage() {
	if t == nil {
		return
	}

	if t.onBeverageStart != nil {
		err := t.onBeverageStart()
		if err != nil {
			panic(err)
		}
	}

	t.b.BoilWater()
	t.b.Brew()
	t.b.PourInCup()

	//子类可以重写该方法来决定是否执行下面动作
	if t.b.WantAddThings() == true {
		t.b.AddThings()
	}
}

// 具体的模板子类 制作咖啡
type MakeCaffee struct {
	template //继承模板
}

func (mc *MakeCaffee) DoOnBeforeBeverage(fn func() error) {
	mc.template.onBeverageStart = fn
}

func NewMakeCaffee() *MakeCaffee {
	makeCaffe := new(MakeCaffee)
	//b 为Beverage，是MakeCaffee的接口，这里需要给接口赋值，指向具体的子类对象
	//来触发b全部接口方法的多态特性。
	makeCaffe.b = makeCaffe
	return makeCaffe
}

func (mc *MakeCaffee) BoilWater() {
	fmt.Println("将水煮到100摄氏度")
}

func (mc *MakeCaffee) Brew() {
	fmt.Println("用水冲咖啡豆")
}

func (mc *MakeCaffee) PourInCup() {
	fmt.Println("将充好的咖啡倒入陶瓷杯中")
}

func (mc *MakeCaffee) AddThings() {
	fmt.Println("添加牛奶和糖")
}

func (mc *MakeCaffee) WantAddThings() bool {
	return true //启动Hook条件
}

// 具体的模板子类 制作茶
type MakeTea struct {
	template //继承模板
}

func (mt *MakeTea) DoOnBeforeBeverage(fn func() error) {
	mt.onBeverageStart = fn
}

func NewMakeTea() *MakeTea {
	makeTea := new(MakeTea)
	//b 为Beverage，是MakeTea，这里需要给接口赋值，指向具体的子类对象
	//来触发b全部接口方法的多态特性。
	makeTea.b = makeTea
	return makeTea
}

func (mt *MakeTea) BoilWater() {
	fmt.Println("将水煮到80摄氏度")
}

func (mt *MakeTea) Brew() {
	fmt.Println("用水冲茶叶")
}

func (mt *MakeTea) PourInCup() {
	fmt.Println("将充好的咖啡倒入茶壶中")
}

func (mt *MakeTea) AddThings() {
	fmt.Println("添加柠檬")
}

func (mt *MakeTea) WantAddThings() bool {
	return false //关闭Hook条件
}

func main() {
	//1. 制作一杯咖啡
	makeCoffee := NewMakeCaffee()
	makeCoffee.DoOnBeforeBeverage(func() error {
		fmt.Println("研磨咖啡豆")
		return nil
	})
	makeCoffee.MakeBeverage() //调用固定模板方法

	fmt.Println("------------")

	//2. 制作茶
	makeTea := NewMakeTea()
	makeTea.MakeBeverage()

}
