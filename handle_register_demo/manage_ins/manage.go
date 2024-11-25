package manage_ins

import (
	"context"
	"errors"
	"fmt"
)

// 享元模式 + 防 git 冲突的设计

type EnumType int64

type InstanceManage interface {
	GetEventID() EnumType
	GetEventName() string
	Handle(ctx context.Context) error
}

var (
	handleManage = make(map[EnumType]InstanceManage)
)

var (
	InstanceHaveNotRegErr = errors.New("实例未注册")
)

// Register 注册实例
func Register(key EnumType, ins InstanceManage) {
	if _, exists := handleManage[key]; exists {
		panic(fmt.Sprintf("Handler with key '%v' already exists", key))
	}
	handleManage[key] = ins
}

// GetInstance 获取实例
func GetInstance(key EnumType) (InstanceManage, error) {
	h, exists := handleManage[key]
	if !exists {
		return nil, InstanceHaveNotRegErr
	}
	return h, nil
}
