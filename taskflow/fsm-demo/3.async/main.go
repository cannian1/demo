package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/looplab/fsm"
)

// 状态机异步任务取消

func main() {
	// 创建状态机
	f := fsm.NewFSM(
		"idle", // 初始状态
		fsm.Events{
			{Name: "start", Src: []string{"idle"}, Dst: "processing"},
			{Name: "complete", Src: []string{"processing"}, Dst: "done"},
		},
		fsm.Callbacks{
			"leave_idle": func(ctx context.Context, e *fsm.Event) {
				fmt.Printf("当前处于 %v 状态，开始异步任务\n", e.FSM.Current())
				// 标记为异步
				e.Async() // 状态维持在现态，等待异步任务完成后手动转换
			},

			"enter_processing": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("进入 processing 状态.")
			},
			"enter_done": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("进入 done 状态.")
			},
		},
	)

	// 打印初始状态
	fmt.Println("Current state:", f.Current())

	// 触发 start 事件
	err := f.Event(context.Background(), "start")
	var asyncError fsm.AsyncError
	ok := errors.As(err, &asyncError)
	if !ok {
		panic(fmt.Sprintf("此处期望 fsm 正在异步执行,返回 'AsyncError',实际上是%v", err))
	}

	// 模拟异步任务被取消
	{
		go func() { // 监听异步任务是否被取消
			<-asyncError.Ctx.Done()
			fmt.Println("外部监听异步任务被取消")
			fmt.Println("当前状态:", f.Current())
		}()

		// 取消异步任务
		asyncError.CancelTransition()
		time.Sleep(20 * time.Millisecond)
	}

	// 如果没有异步标记，这里才会有 err 触发
	if err = f.Transition(); err != nil {
		fmt.Println("任务转化到下一阶段失败:", err)
	}

	if f.Can("complete") {
		if err = f.Event(context.Background(), "complete"); err != nil {
			fmt.Println("触发 complete 行为时出错:", err)
		}
	} else {
		fmt.Println("当前状态无法触发 complete 行为")
	}

	fmt.Println("最终状态:", f.Current())
}
