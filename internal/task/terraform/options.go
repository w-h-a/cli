package terraform

import (
	"context"

	"github.com/w-h-a/cli/internal/task"
)

func TerraformWithVars(vars map[string]string) task.TaskOption {
	return func(o *task.TaskOptions) {
		o.Context = context.WithValue(o.Context, "tf_vars_key", vars)
	}
}

func TerraformWithRemoteStates(rs map[string]string) task.TaskOption {
	return func(o *task.TaskOptions) {
		o.Context = context.WithValue(o.Context, "tf_remote_states_key", rs)
	}
}
