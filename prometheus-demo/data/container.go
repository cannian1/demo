package data

import (
	"context"
	"demo/prometheus-demo/executor"
)

type ContainerHandler interface {
	// GetDiskUsage 获取磁盘使用情况
	GetDiskUsage(ctx context.Context) (string, error)
}

type Container struct {
	container string // 容器 ID 或 容器名称
}

func NewContainer(container string) *Container {
	return &Container{container: container}
}

func (c Container) GetDiskUsage(ctx context.Context) (string, error) {
	exec := executor.NewCmdExecutor()

	dockerArgs := []string{"exec", c.container, "df", "-h"}
	output, err := exec.Run(ctx, "docker", dockerArgs...)
	if err != nil {
		return "", err
	}
	return output, nil
}
