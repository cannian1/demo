package main

import "fmt"

// Implementor 接口（定义实现部分的接口）
type DrawingAPI interface {
	DrawCircle(x, y, radius int)
}

// 具体实现：Windows 绘制
type WindowsDrawingAPI struct{}

func (w *WindowsDrawingAPI) DrawCircle(x, y, radius int) {
	fmt.Printf("Windows: 绘制圆形 [中心点: (%d, %d), 半径: %d]\n", x, y, radius)
}

// 具体实现：Linux 绘制
type LinuxDrawingAPI struct{}

func (l *LinuxDrawingAPI) DrawCircle(x, y, radius int) {
	fmt.Printf("Linux: 绘制圆形 [中心点: (%d, %d), 半径: %d]\n", x, y, radius)
}

// Abstraction 抽象类（定义形状接口）
type Shape interface {
	Draw()
}

// Refined Abstraction：圆形
type Circle struct {
	x, y, radius int
	drawingAPI   DrawingAPI // 持有实现部分的接口
}

// 构造函数
func NewCircle(x, y, radius int, drawingAPI DrawingAPI) *Circle {
	return &Circle{x: x, y: y, radius: radius, drawingAPI: drawingAPI}
}

// 使用具体实现绘制圆形
func (c *Circle) Draw() {
	c.drawingAPI.DrawCircle(c.x, c.y, c.radius)
}

func main() {
	// 使用 Windows 实现绘制圆形
	windowsCircle := NewCircle(10, 20, 30, &WindowsDrawingAPI{})
	windowsCircle.Draw()

	// 使用 Linux 实现绘制圆形
	linuxCircle := NewCircle(50, 60, 70, &LinuxDrawingAPI{})
	linuxCircle.Draw()
}
