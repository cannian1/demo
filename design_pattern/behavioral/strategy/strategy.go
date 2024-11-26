package main

import "fmt"

// 定义策略接口
type PaymentStrategy interface {
	Pay(amount float64) error
}

// 实现具体的支付策略
type CreditCardStrategy struct {
	name       string
	cardNumber string
	date       string
}

func (c *CreditCardStrategy) Pay(amount float64) error {
	fmt.Printf("Paying %0.2f using credit card\n", amount)
	return nil
}

type PayPalStrategy struct {
	email    string
	password string
}

func (p *PayPalStrategy) Pay(amount float64) error {
	fmt.Printf("Paying %0.2f using PayPal\n", amount)
	return nil
}

// 定义上下文类
type PaymentContext struct {
	amount   float64
	strategy PaymentStrategy
}

func NewPaymentContext(amount float64, strategy PaymentStrategy) *PaymentContext {
	return &PaymentContext{
		amount:   amount,
		strategy: strategy,
	}
}

func (p *PaymentContext) Pay() error {
	return p.strategy.Pay(p.amount)
}

func main() {
	creditCardStrategy := &CreditCardStrategy{
		name:       "John Doe",
		cardNumber: "1234 5678 9012 3456",
		date:       "01/22",
	}
	paymentContext := NewPaymentContext(100, creditCardStrategy)
	paymentContext.Pay()

	payPalStrategy := &PayPalStrategy{
		email:    "john.doe@example.com",
		password: "114514",
	}
	paymentContext = NewPaymentContext(200, payPalStrategy)
	paymentContext.Pay()
}
