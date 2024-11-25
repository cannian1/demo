package main

import "fmt"

type Car struct {
	Brand string
	Price float64
}

type PriceDecorator interface {
	DecoratePrice(car Car) Car
}

type ExtraPriceDecorator struct {
	ExtraPrice float64
}

func (d ExtraPriceDecorator) DecoratePrice(car Car) Car {
	car.Price += d.ExtraPrice
	return car
}

type DiscountPriceDecorator struct {
	Percent float64
}

func (d DiscountPriceDecorator) DecoratePrice(car Car) Car {
	car.Price = car.Price * d.Percent * 0.01
	return car
}

func main() {
	toyota := Car{Brand: "Toyota", Price: 10000}
	decorator := ExtraPriceDecorator{ExtraPrice: 500}
	decoratedCar := decorator.DecoratePrice(toyota)
	fmt.Printf("%+v\n", decoratedCar)

	byd := Car{Brand: "byd", Price: 10000}
	decorator2 := DiscountPriceDecorator{Percent: 10}
	decoratedCar2 := decorator2.DecoratePrice(byd)
	fmt.Printf("%+v\n", decoratedCar2)
}
