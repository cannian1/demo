package domain

import (
	"demo/handle_register_demo/domain/ins_b"
	"demo/handle_register_demo/manage_ins"
)

func init() {
	bbb := ins_b.NewBBB()
	manage_ins.Register(ins_b.BBBType, bbb)
}
