package main

import "fmt"

// 定义多个子系统和组件

type CPU struct {
	speed float64
}

func (g *CPU) overLock() {
	g.speed *= 10
}

func (c *CPU) run() {
	fmt.Println("CPU running at", c.speed, "GHz")
}

type Memory struct {
	capacity int
}

func (m *Memory) load() {
	fmt.Println("Memory loaded with", m.capacity, "GB")
}

type Disk struct {
	size int
}

func (d *Disk) read() {
	fmt.Println("Disk reading data with size", d.size, "KB/s")
}

type GPU struct {
	mode string
}

func (g *GPU) overLock() {
	g.mode = "独显输出"
}

func (g *GPU) run() {
	fmt.Printf("gpu run in %s mode\n", g.mode)
}

// 定义外观接口

type ComputerFacade struct {
	cpu    *CPU
	memory *Memory
	disk   *Disk
	gpu    *GPU
}

func NewComputerFacade(speed float64, capacity int, size int) *ComputerFacade {
	cpu := &CPU{speed: speed}
	memory := &Memory{capacity: capacity}
	disk := &Disk{size: size}
	gpu := &GPU{mode: "混合输出"}

	return &ComputerFacade{cpu: cpu, memory: memory, disk: disk, gpu: gpu}
}

func (f *ComputerFacade) Run() {
	fmt.Println("Computer starting...")
	f.cpu.run()
	f.memory.load()
	f.disk.read()
	fmt.Println("Computer running...")
}

func (f *ComputerFacade) OverclockingRun() {
	fmt.Println("Computer starting...")
	f.cpu.overLock()
	f.gpu.overLock()

	f.cpu.run()
	f.memory.load()
	f.disk.read()
	f.gpu.run()
	fmt.Println("Computer running...")
}

func main() {
	computer := NewComputerFacade(3.4, 16, 1024)
	computer.Run()

	fmt.Println("---------")

	computer.OverclockingRun()
}
