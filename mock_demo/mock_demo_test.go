package mock_demo

import (
	"fmt"
	"testing"

	"github.com/prashantv/gostub"
	"go.uber.org/mock/gomock"
)

// 自定义匹配器，匹配大于指定值的参数
type GreaterThanMatcher struct {
	Threshold string
}

// 实现 Matches 方法，匹配是否大于阈值
func (m GreaterThanMatcher) Matches(x any) bool {
	if value, ok := x.(string); ok {
		return value > m.Threshold
	}
	return false
}

// 实现 String 方法，返回匹配器的描述
func (m GreaterThanMatcher) String() string {
	return fmt.Sprintf("比 %s 大", m.Threshold)
}

// 创建自定义匹配器
func GreaterThan(threshold string) gomock.Matcher {
	return GreaterThanMatcher{Threshold: threshold}
}

func TestMC_WriteAndSend(t *testing.T) {
	// 和 stub 一样，存根函数，执行完测试恢复现场
	//old := getSign
	//getSign = func() string {
	//	return ""
	//}
	//defer func() {
	//	getSign = old
	//}()

	stubs := gostub.Stub(&getSign, func() string {
		return "2024年10月13日 23:02:00"
	})
	defer stubs.Reset()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMail := NewMockMail(mockCtrl)

	// 对 Mail 接口的 sendMail 方法进行了 mock
	// 对于传参的 mock 函数期望校验有以下几个方法
	// - gomock.Eq() 是否相等
	// - gomock.AnyOf() 其中之一
	// - gomock.InAnyOrder() 任意顺序
	// - gomock.Len() 参数长度
	// - gomock.AssignableToTypeOf() 参数类型，里面随便放一个和期望类型相同的变量就可以
	// - gomock.Not() 不为
	// - gomock.Any() 任意
	// 还可以自定义匹配器，需要实现 Matches 和 String 方法，如上所示
	mockMail.EXPECT().
		sendMail(gomock.Eq("x1"), GreaterThan("x0"), gomock.Any(), "某人于2024年10月13日 23:02:00"). // mock 函数期望收到的传参
		Return(nil).                                                                           // mock 函数被调用后期望的返回值
		Times(1)                                                                               // 调用预期执行的次数

	// 用 mock 的方法（已经实现了 Mail 接口）实例化 MC
	mc := NewMC(mockMail)

	// 尝试调用 WriteAndSend，请求的参数
	mc.WriteAndSend("x1", "x2", "x3", "某人于")
}
