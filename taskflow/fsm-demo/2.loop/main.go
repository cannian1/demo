package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/looplab/fsm"
)

func main() {
	// 定义有限状态机
	processFSM := fsm.NewFSM(
		"empty", // 初始状态为 "empty"
		fsm.Events{
			// 从 "empty" 到 "ready"
			{Name: "fork", Src: []string{"empty"}, Dst: "ready"},
			// 从 "ready" 到 "running"
			{Name: "dispatch", Src: []string{"ready"}, Dst: "running"},
			// 从 "running" 到 "running" (继续工作)
			{Name: "work", Src: []string{"running"}, Dst: "running"},
			// 从 "running" 到 "ready" (时间片用完)
			{Name: "time_slice", Src: []string{"running"}, Dst: "ready"},
			// 从 "running" 到 "blocked" (等待资源)
			{Name: "wait", Src: []string{"running"}, Dst: "blocked"},
			// 从 "blocked" 到 "blocked" (继续等待资源)
			{Name: "wait", Src: []string{"blocked"}, Dst: "blocked"},
			// 从 "blocked" 到 "ready" (资源就绪)
			{Name: "resume", Src: []string{"blocked"}, Dst: "ready"},
			// 从 "ready" 到 "ready" (继续等待)
			{Name: "wait", Src: []string{"ready"}, Dst: "ready"},
			// 从 "running" 到 "finished" (进程终止)
			{Name: "release", Src: []string{"running"}, Dst: "finished"},
		},
		fsm.Callbacks{
			"dispatch": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Dispatch: Process is now running.")
			},
			"time_slice": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Time Slice: Process is now ready.")
			},
			"wait": func(_ context.Context, e *fsm.Event) {
				if e.Src == "running" {
					fmt.Println("Wait: Process is now blocked.")
				} else {
					fmt.Println("Wait: Process is now ready.")
				}
			},
			"resume": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Resume: Process is now ready.")
			},
			"work": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Work: Process is still running.")
			},
			"kill": func(_ context.Context, e *fsm.Event) {
				killedErr := errors.New("进程无法被终止")
				checkStatus := func(pid int) error {
					if pid == 1 && e.Src == "running" {
						return killedErr
					}
					return nil
				}

				pid := e.Args[0].(map[string]int)["pid"]
				if err := checkStatus(pid); err != nil {
					fmt.Println("Kill: Process is not killed.")
					// 取消状态转换并返回错误
					e.FSM.SetState(e.Src)
					e.Cancel(err)
					return
				} else {
					fmt.Println("Kill: Process is killed.")
				}
			},
		},
	)

	// 打印初始状态
	fmt.Println("初始状态:", processFSM.Current())

	// 模拟状态转换
	err := processFSM.Event(context.Background(), "fork") // empty -> ready
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	err = processFSM.Event(context.Background(), "dispatch") // ready -> running
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	err = processFSM.Event(context.Background(), "work") // running -> running
	if err != nil && !errors.Is(err, fsm.NoTransitionError{}) {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	err = processFSM.Event(context.Background(), "time_slice") // running -> ready
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	err = processFSM.Event(context.Background(), "dispatch") // ready -> running
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	err = processFSM.Event(context.Background(), "wait") // running -> blocked
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	err = processFSM.Event(context.Background(), "resume") // blocked -> ready
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	processFSM.SetState("running")
	err = processFSM.Event(context.Background(), "release", map[string]int{"pid": 1}) // running -> finished
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("当前状态:", processFSM.Current())

	// fmt.Println(fsm.Visualize(processFSM))
}
