package domain

import (
	"demo/handle_register_demo/domain/ins_a"
	"demo/handle_register_demo/manage_ins"
)

func init() {
	aaa := ins_a.NewAAA()
	manage_ins.Register(ins_a.AAAType, aaa)
}
