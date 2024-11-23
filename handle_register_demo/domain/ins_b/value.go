package ins_b

import (
	"context"
	"demo/handle_register_demo/manage_ins"
	"fmt"
)

var (
	BBBType manage_ins.EnumType = 2
	BBBName string              = "BBB"
)

func NewBBB() BBB {
	return BBB{}
}

type BBB struct {
}

func (bbb BBB) GetEventID() manage_ins.EnumType {
	return BBBType
}

func (bbb BBB) GetEventName() string {
	return BBBName
}

func (bbb BBB) Handle(ctx context.Context) error {
	fmt.Println("BBB 执行中")
	return nil
}
