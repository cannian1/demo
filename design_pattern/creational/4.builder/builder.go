package main

import "fmt"

// 产品
type Character struct {
	profession string
	weapon     string
	equipment  string
}

// 抽象建造者
type CharacterBuilder interface {
	SetProfession()
	SetWeapon()
	SetEquipment()
	GetCharacter() Character
}

// 具体建造者：战士，用于创建战士对象
type WarriorBuilder struct {
	character Character
}

func (wb *WarriorBuilder) SetProfession() {
	wb.character.profession = "战士"
}

func (wb *WarriorBuilder) SetWeapon() {
	wb.character.weapon = "剑"
}

func (wb *WarriorBuilder) SetEquipment() {
	wb.character.equipment = "盾牌"
}

func (wb *WarriorBuilder) GetCharacter() Character {
	return wb.character
}

// 具体建造者：法师，用于创建法师对象
type MageBuilder struct {
	character Character
}

func (mb *MageBuilder) SetProfession() {
	mb.character.profession = "法师"
}

func (mb *MageBuilder) SetWeapon() {
	mb.character.weapon = "法杖"
}

func (mb *MageBuilder) SetEquipment() {
	mb.character.equipment = "魔导书"
}

func (mb *MageBuilder) GetCharacter() Character {
	return mb.character
}

// 指挥者：控制建造的流程
type Director struct {
	builder CharacterBuilder
}

func (d *Director) Construct() Character {
	d.builder.SetProfession()
	d.builder.SetWeapon()
	d.builder.SetEquipment()
	return d.builder.GetCharacter()
}

// 示例代码
func main() {

	director := Director{}

	// 构建战士角色
	warriorBuilder := &WarriorBuilder{}
	director.builder = warriorBuilder
	warrior := director.Construct()
	fmt.Println("战士角色：", warrior)

	// 构建法师角色
	mageBuilder := &MageBuilder{}
	director.builder = mageBuilder
	mage := director.Construct()
	fmt.Println("法师角色：", mage)
}
