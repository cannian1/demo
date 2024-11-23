package ins_a

import (
	"context"
	"demo/handle_register_demo/manage_ins"
	"fmt"
)

var (
	AAAType manage_ins.EnumType = 1
	AAAName string              = "AAA"
)

func NewAAA() AAA {
	return AAA{}
}

type AAA struct {
}

func (aaa AAA) GetEventID() manage_ins.EnumType {
	return AAAType
}

func (aaa AAA) GetEventName() string {
	return AAAName
}

func (aaa AAA) Handle(ctx context.Context) error {
	fmt.Println("AAA 执行中")
	return nil
}
