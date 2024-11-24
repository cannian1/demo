package main

import (
	"fmt"
	"log/slog"
)

type Logger interface {
	Info(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
	Warn(v ...interface{})
}

var _ Logger = &DefaultLog{}

type DefaultLog struct {
}

func (d DefaultLog) Info(v ...interface{}) {
	slog.Info("Msg", "Ext:", fmt.Sprint(v...))
}

func (d DefaultLog) Error(v ...interface{}) {
	slog.Error("Msg", "Ext:", fmt.Sprint(v...))
}

func (d DefaultLog) Debug(v ...interface{}) {
	slog.Debug("Msg", "Ext:", fmt.Sprint(v...))
}

func (d DefaultLog) Warn(v ...interface{}) {
	slog.Warn("Msg", "Ext:", fmt.Sprint(v...))
}

// ThirdPartyLogger 是模拟的第三方日志库，具有不同的方法签名
type ThirdPartyLogger struct{}

func (t *ThirdPartyLogger) PrintInfo(msg string) {
	fmt.Println("[INFO]:", msg)
}

func (t *ThirdPartyLogger) PrintError(msg string) {
	fmt.Println("[ERROR]:", msg)
}

func (t *ThirdPartyLogger) PrintDebug(msg string) {
	fmt.Println("[DEBUG]:", msg)
}

func (t *ThirdPartyLogger) PrintWarn(msg string) {
	fmt.Println("[WARN]:", msg)
}

// Adapter 将 ThirdPartyLogger 适配到 Logger 接口
type Adapter struct {
	thirdPartyLogger *ThirdPartyLogger
}

// 实现 Logger 接口的方法，将调用委托给 ThirdPartyLogger
func (a *Adapter) Info(v ...interface{}) {
	a.thirdPartyLogger.PrintInfo(fmt.Sprint(v...))
}

func (a *Adapter) Error(v ...interface{}) {
	a.thirdPartyLogger.PrintError(fmt.Sprint(v...))
}

func (a *Adapter) Debug(v ...interface{}) {
	a.thirdPartyLogger.PrintDebug(fmt.Sprint(v...))
}

func (a *Adapter) Warn(v ...interface{}) {
	a.thirdPartyLogger.PrintWarn(fmt.Sprint(v...))
}

func main() {

	// logger := &DefaultLog{}

	// 创建第三方日志实例
	thirdPartyLogger := &ThirdPartyLogger{}

	// 创建适配器实例，将第三方日志传入适配器
	logger := &Adapter{
		thirdPartyLogger: thirdPartyLogger,
	}

	// 使用 Logger 接口记录日志
	logger.Info("This is an info message")
	logger.Error("This is an error message")
	logger.Debug("This is a debug message")
	logger.Warn("This is a warning message")
}
