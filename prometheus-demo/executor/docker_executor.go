package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerExecutor 指定容器执行命令，但是返回结果有乱码，不如标准库那种方式获取
type DockerExecutor struct {
	cli       *client.Client
	container string
}

func NewDockerExecutor(cli *client.Client, container string) *DockerExecutor {
	return &DockerExecutor{cli: cli, container: container}
}

func (d *DockerExecutor) Run(ctx context.Context, command string, args ...string) (string, error) {
	dockerCmd := append([]string{command}, args...)

	resp, err := d.cli.ContainerExecCreate(ctx, d.container, container.ExecOptions{
		Cmd:          dockerCmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %w", err)
	}

	// 附加到 exec 会话并获取输出
	execID := resp.ID
	respAttach, err := d.cli.ContainerExecAttach(ctx, execID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec: %v", err)
	}
	defer respAttach.Close()

	// 读取输出
	var output bytes.Buffer
	_, err = io.Copy(&output, respAttach.Reader)
	return output.String(), err
}

func (d *DockerExecutor) RunWithEnv(ctx context.Context, command string, args []string, env map[string]string) (string, error) {
	panic("implement me")
}
