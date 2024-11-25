package main

import "fmt"

// Component 接口定义了组合对象和叶子对象的公共方法。
type Component interface {
	GetName() string
	List(indent int)
}

// File 是叶子对象，它代表文件。
type File struct {
	name string
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) List(indent int) {
	fmt.Printf("%s- %s\n", getIndent(indent), f.name)
}

// Directory 是组合对象，它代表文件夹。
type Directory struct {
	name     string
	children []Component
}

func (d *Directory) GetName() string {
	return d.name
}

func (d *Directory) List(indent int) {
	fmt.Printf("%s+ %s\n", getIndent(indent), d.name)
	for _, child := range d.children {
		child.List(indent + 2)
	}
}

// AddChild 方法用于向文件夹中添加子组件。
func (d *Directory) AddChild(c Component) *Directory {
	d.children = append(d.children, c)

	return d
}

// RemoveChild 方法用于从文件夹中移除子组件。
func (d *Directory) RemoveChild(c Component) *Directory {
	for i, child := range d.children {
		if child == c {
			d.children = append(d.children[:i], d.children[i+1:]...)
			return d
		}
	}

	return d
}

// getIndent 方法返回指定缩进级别的空格字符串。
func getIndent(indent int) string {
	result := ""
	for i := 0; i < indent; i++ {
		result += " "
	}
	return result
}

func main() {
	// 创建文件系统
	root := &Directory{name: "root"}
	documents := &Directory{name: "documents"}
	pictures := &Directory{name: "pictures"}
	music := &Directory{name: "music"}
	file1 := &File{name: "file1.txt"}
	file2 := &File{name: "file2.txt"}
	file3 := &File{name: "file3.txt"}
	popMusic := &Directory{name: "pop music"}

	root.AddChild(documents).
		AddChild(pictures).
		AddChild(music)
	documents.
		AddChild(file1).
		AddChild(file2)
	pictures.AddChild(file3)
	music.AddChild(popMusic)

	// 打印文件系统的层次结构
	root.List(0)
}
