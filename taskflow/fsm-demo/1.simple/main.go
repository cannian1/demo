package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/looplab/fsm"
)

// FontName 字体名称，来自操作系统
type FontName string

const (
	// MicrosoftYaHei 微软雅黑
	MicrosoftYaHei FontName = "Microsoft YaHei"
	// SimSun 宋体
	SimSun FontName = "SimSun"
)

func main() {

	f := fsm.NewFSM(
		"无订单", // 初始状态
		fsm.Events{
			{Name: "新建", Src: []string{"无订单"}, Dst: "订单被创建"},
			{Name: "支付", Src: []string{"订单被创建"}, Dst: "订单已支付"},
			{Name: "发货", Src: []string{"订单已支付"}, Dst: "订单已发货"},
			{Name: "收货", Src: []string{"订单已发货"}, Dst: "订单已收货"},
			{Name: "完成", Src: []string{"订单已收货"}, Dst: "订单已完成"},
			{Name: "冻结", Src: []string{"订单被创建", "订单已支付"}, Dst: "订单已冻结"},
		},
		fsm.Callbacks{
			// 在进入状态时触发的回调函数
			"enter_state": func(ctx context.Context, event *fsm.Event) {
				fmt.Println("【记录日志】进入目标状态：", event.Dst)
			},
			// 在离开状态时触发的回调函数
			"leave_state": func(ctx context.Context, event *fsm.Event) {
				fmt.Println("【记录日志】离开先前状态：", event.Src)
			},
			"enter_订单被创建": func(ctx context.Context, event *fsm.Event) {
				orderID := uuid.New().String()
				// 设置元数据
				event.FSM.SetMetadata("orderID", orderID)
				fmt.Printf("【记录日志】订单 %v 被创建\n", orderID)
			},
			"enter_订单已支付": func(ctx context.Context, event *fsm.Event) {
				// FSM 的 Event 方法的可变长参数会传递给事件回调函数
				if len(event.Args) == 0 {
					return
				}
				extraInfo := event.Args[0]
				fmt.Println("【进入订单已支付事件回调，获取客户备注】", extraInfo.(string))
			},
			"enter_订单已冻结": func(ctx context.Context, event *fsm.Event) {
				// 根据事件来源判断是否记录日志
				orderID, _ := event.FSM.Metadata("orderID")

				if event.Src == "订单已支付" {
					fmt.Printf("【记录日志】订单 %s 已支付，冻结订单\n", orderID)
				} else {
					fmt.Println("【记录日志】订单被创建，冻结订单", orderID)
				}
			},
		},
	)

	// 可视化状态机，生成 DOT 文件和图片
	//createDotFile(f)
	//createImageFile()

	fmt.Println("状态机初始状态为：", f.Current())

	// 触发【新建】事件
	err := f.Event(context.Background(), "新建")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("触发【新建】行为后，状态机的状态变为：", f.Current())

	// 触发【支付】事件
	// Event 后面的可变长参数会传递给 Callbacks 中的事件回调函数
	err = f.Event(context.Background(), "支付", "老板，我今天生日可以送我一个小菜么")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("触发【支付】行为后，状态机的状态变为：", f.Current())

	// -------------------
	//// 触发冻结订单事件
	//err = f.Event(context.Background(), "冻结")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	// -------------------

	// 触发【发货】事件
	err = f.Event(context.Background(), "发货")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("触发【发货】行为后，状态机的状态变为：", f.Current())

	// 触发【收货】事件
	err = f.Event(context.Background(), "收货")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("触发【收货】行为后，状态机的状态变为：", f.Current())

	// 触发【完成】事件
	err = f.Event(context.Background(), "完成")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("触发【完成】行为后，状态机的状态变为：", f.Current())
	// 获取元数据
	orderID, ok := f.Metadata("orderID")
	if ok {
		fmt.Printf("订单ID：%v 已完成\n", orderID)
		return
	}
}

// 生成状态机的 DOT 格式文件
func createDotFile(f *fsm.FSM) error {
	dotStr := adaptToChinese(fsm.Visualize(f), SimSun)
	fmt.Println(dotStr)

	if err := os.WriteFile("fsm.dot", []byte(dotStr), 0644); err != nil {
		fmt.Println("生成状态机 DOT 文件失败：", err)
		return err
	}
	return nil
}

// 调用 graphviz 将状态机的 DOT 格式文件转换为图片
func createImageFile() error {
	cmd := exec.Command("dot", "-Tpng", "fsm.dot", "-o", "fsm.png")
	err := cmd.Run()
	if err != nil {
		fmt.Println("生成状态机图片失败：", err)
		return err
	}
	return nil
}

// 适配中文：将 Graphviz 的 DOT 文件中的节点和边字体设置为中文，在第一行后追加
// node [fontname="SimSun"];
// edge [fontname="SimSun"];
func adaptToChinese(input string, fontName FontName) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var lines []string

	// 遍历每一行
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if len(lines) == 1 { // 如果是第一行，追加字符串
			str := fmt.Sprintf("\tnode [fontname=\"%s\"];\n\tedge [fontname=\"%s\"];", fontName, fontName)
			lines = append(lines, str)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}

	// 将所有行合并为一个字符串
	return strings.Join(lines, "\n")
}
