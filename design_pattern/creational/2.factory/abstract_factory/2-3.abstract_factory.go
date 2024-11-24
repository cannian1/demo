package main

import "fmt"

// ======= 抽象层 =========

type AbstractGPU interface {
	Display()
}

type AbstractCPU interface {
	Calculate()
}

type AbstractMemory interface {
	Storage()
}

type AbstractComputerFactory interface {
	CreateGPU() AbstractGPU
	CreateCPU() AbstractCPU
	CreateMemory() AbstractMemory
}

// ======== 实现层 =========

type NvidiaCPU struct{}

func (cpu *NvidiaCPU) Calculate() {
	fmt.Println("CPU is NvidiaCPU ")
}

type NvidiaGPU struct{}

func (gpu *NvidiaGPU) Display() {
	fmt.Println("GPU is NvidiaGPU ")
}

type NvidiaMemory struct{}

func (m *NvidiaMemory) Storage() {
	fmt.Println("Memory is NvidiaMemory ")
}

type IntelCPU struct{}

func (cpu *IntelCPU) Calculate() {
	fmt.Println("CPU is IntelCPU")
}

type IntelGPU struct{}

func (gpu *IntelGPU) Display() {
	fmt.Println("GPU is IntelGPU")
}

type IntelMemory struct{}

func (m *IntelMemory) Storage() {
	fmt.Println("Memory is IntelMemory ")
}

type KingstonCPU struct{}

func (cpu *KingstonCPU) Calculate() {
	fmt.Println("CPU is KingstonCPU ")
}

type KingstonGPU struct{}

func (gpu *KingstonGPU) Display() {
	fmt.Println("GPU is  KingstonGPU")
}

type KingstonMemory struct{}

func (m *KingstonMemory) Storage() {
	fmt.Println("Memory is KingstonMemory ")
}

// Intel工厂
type IntelFactory struct{}

func (factory *IntelFactory) CreateGPU() AbstractGPU {
	var product AbstractGPU
	product = new(IntelGPU)
	return product
}
func (factory *IntelFactory) CreateCPU() AbstractCPU {
	var product AbstractCPU
	product = new(IntelCPU)
	return product
}
func (factory *IntelFactory) CreateMemory() AbstractMemory {
	var product AbstractMemory
	product = new(IntelMemory)
	return product
}

// Nvidia工厂
type NvidiaFactory struct{}

func (factory *NvidiaFactory) CreateGPU() AbstractGPU {
	var product AbstractGPU
	product = new(NvidiaGPU)
	return product
}
func (factory *NvidiaFactory) CreateCPU() AbstractCPU {
	var product AbstractCPU
	product = new(NvidiaCPU)
	return product
}
func (factory *NvidiaFactory) CreateMemory() AbstractMemory {
	var product AbstractMemory
	product = new(NvidiaMemory)
	return product
}

// Kingston工厂
type KingstonFactory struct{}

func (factory *KingstonFactory) CreateGPU() AbstractGPU {
	var product AbstractGPU
	product = new(KingstonGPU)
	return product
}
func (factory *KingstonFactory) CreateCPU() AbstractCPU {
	var product AbstractCPU
	product = new(KingstonCPU)
	return product
}
func (factory *KingstonFactory) CreateMemory() AbstractMemory {
	var product AbstractMemory
	product = new(KingstonMemory)
	return product
}

// ======== 业务逻辑层 =======
type Computer struct {
	CPU    AbstractCPU
	GPU    AbstractGPU
	Memory AbstractMemory
}

func (c *Computer) Stats() {
	fmt.Println("-------------电脑配置单---------------")
	c.CPU.Calculate()
	c.GPU.Display()
	c.Memory.Storage()
}

func main() {

	intelFactory := new(IntelFactory)
	nvidiaFactor := new(NvidiaFactory)
	kingstonFactory := new(KingstonFactory)

	computer1 := Computer{
		GPU:    intelFactory.CreateGPU(),
		CPU:    intelFactory.CreateCPU(),
		Memory: intelFactory.CreateMemory(),
	}
	computer2 := Computer{
		GPU:    nvidiaFactor.CreateGPU(),
		CPU:    intelFactory.CreateCPU(),
		Memory: kingstonFactory.CreateMemory(),
	}
	computer1.Stats()
	computer2.Stats()

}
