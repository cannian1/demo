package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type CmdExecutor struct{}

func NewCmdExecutor() *CmdExecutor {
	return &CmdExecutor{}
}

// Run 执行 command 携带 args 参数，返回输出
func (e *CmdExecutor) Run(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if stderr.Len() > 0 {
		return "", errors.New(stderr.String())
	}
	if err != nil {
		return "", err
	}

	return stdout.String(), err
}

// RunWithEnv 执行 command 携带 args 参数，自定义环境变量
func (e *CmdExecutor) RunWithEnv(ctx context.Context, command string, args []string, env map[string]string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if env != nil {
		envVars := os.Environ()
		for key, value := range env {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = envVars
	}

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}
