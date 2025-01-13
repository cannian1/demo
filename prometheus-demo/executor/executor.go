package executor

import "context"

type Executor interface {
	Run(ctx context.Context, command string, args ...string) (string, error)
	RunWithEnv(ctx context.Context, command string, args []string, env map[string]string) (string, error)
}
