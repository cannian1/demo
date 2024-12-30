package bind_public_IP

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/looplab/fsm"
)

// 领域是裸金属主机绑定公网 IP，这里将实际逻辑简化很多
// 原业务支持多台主机绑定，并且根据主机的归属情况选择 SNAT 和 DNAT 策略

const (
	initialState         = "初始状态"
	unexpectParam        = "被选中的机器存在不被期望的状态，可能存在脏数据"
	snatPolicyIsSet      = "SNAT策略已配置"
	dnatPolicyIsSet      = "DNAT策略已配置"
	bandwidthPolicyIsSet = "带宽策略已配置"
	bindStatusAreChanged = "绑定状态修改成功"
	bindPublicIPFinish   = "正常绑定完成，流程结束"
	needUserSelfFinish   = "剩余配置需要用户自行分配，流程结束"
	unexpectFinish       = "不符合预期的失败出现，流程结束"
)

const (
	setSNatPolicy      = "配置SNAT策略"
	setDNatPolicy      = "配置DNAT策略"
	setBandwidthPolicy = "配置带宽策略"
	changeBindStatus   = "修改绑定状态"
	completeBinding    = "完成绑定"
	letUserConfigIt    = "让用户自行配置"

	triggerUnexpectParam = "校验到不满足所有配置所需的依赖"
	thirdPartyError      = "第三方服务执行异常"
)

type BindPublicIP struct {
	TaskID     int
	MachineIDs []int  // 机器 ID
	Uid        uint64 // 要绑定的用户 ID
	ZoneID     int    // 可用区 ID
	FwID       int    // 防火墙 ID

	//SNATPortType int // SNAT 只有一种配置方式
	DNATPortType int // DNAT 端口类型 0：全端口映射 1：部分端口映射

	PublicIP      string // 公网 IP
	PrivateIP     string // 内网 IP
	BandwidthUp   int    // 带宽（上行）
	BandwidthDown int    // 带宽（下行）
	Status        string // 状态

	FSM *fsm.FSM
}

func NewBindPublicIP(taskID int, machineIDs []int, uid uint64, zoneID int, fwID int, DNATPortType int, publicIP string, privateIP string, bandwidthUp int, bandwidthDown int, status string, FSM *fsm.FSM) *BindPublicIP {
	return &BindPublicIP{
		TaskID:        taskID,
		MachineIDs:    machineIDs,
		Uid:           uid,
		ZoneID:        zoneID,
		FwID:          fwID,
		DNATPortType:  DNATPortType,
		PublicIP:      publicIP,
		PrivateIP:     privateIP,
		BandwidthUp:   bandwidthUp,
		BandwidthDown: bandwidthDown,
		Status:        status,

		FSM: fsm.NewFSM(initialState,
			fsm.Events{
				{Name: setSNatPolicy, Src: []string{initialState}, Dst: snatPolicyIsSet},
				{Name: setDNatPolicy, Src: []string{snatPolicyIsSet}, Dst: dnatPolicyIsSet},
				{Name: setBandwidthPolicy, Src: []string{dnatPolicyIsSet}, Dst: bandwidthPolicyIsSet},
				{Name: changeBindStatus, Src: []string{bandwidthPolicyIsSet}, Dst: bindStatusAreChanged},
				{Name: completeBinding, Src: []string{bindStatusAreChanged}, Dst: bindPublicIPFinish},
				{Name: letUserConfigIt, Src: []string{snatPolicyIsSet}, Dst: needUserSelfFinish},

				// 参数不满足绑定状态流转的情况：
				// 初始状态：存在已绑定主机
				{Name: triggerUnexpectParam, Src: []string{initialState, dnatPolicyIsSet}, Dst: unexpectParam},
				{Name: thirdPartyError, Src: []string{snatPolicyIsSet, dnatPolicyIsSet, bandwidthPolicyIsSet, bindStatusAreChanged}, Dst: unexpectFinish},
			},
			fsm.Callbacks{
				"enter_state": func(ctx context.Context, event *fsm.Event) {
					fmt.Printf("进入状态：%s\n", event.Dst)
				},
			},
		),
	}
}

// SetSNAT 设置 SNAT 规则
func (b *BindPublicIP) SetSNAT(ctx context.Context) error {
	// 0. 校验参数
	// 检查是否存在已绑定主机
	fmt.Println("校验参数, 检查是否存在已绑定主机")
	// 如果存在就流转到 unexpectParam
	// b.FSM.Event(ctx, triggerUnexpectParam)

	// 1. 调用第三方服务配置 SNAT 规则
	fmt.Println("配置 SNAT 规则")

	// 检查是否配置成功 rpc

	// 配置成功后，流转到配置 DNAT 规则
	return b.FSM.Event(ctx, setSNatPolicy)
}

func (b *BindPublicIP) Visualize() string {
	return fsm.Visualize(b.FSM)
}

// 生成状态机的 DOT 格式文件
func createDotFile(f *fsm.FSM) error {
	dotStr := adaptToChinese(fsm.Visualize(f), "SimSun")
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
func adaptToChinese(input string, fontName string) string {
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
