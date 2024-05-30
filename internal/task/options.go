package task

import "context"

type TaskOption func(o *TaskOptions)

type TaskOptions struct {
	Name    string
	Source  string
	Path    string
	EnvVars map[string]string
	Context context.Context
}

func TaskWithName(n string) TaskOption {
	return func(o *TaskOptions) {
		o.Name = n
	}
}

func TaskWithSource(s string) TaskOption {
	return func(o *TaskOptions) {
		o.Source = s
	}
}

func TaskWithPath(p string) TaskOption {
	return func(o *TaskOptions) {
		o.Path = p
	}
}

func TaskWithEnvVars(ev map[string]string) TaskOption {
	return func(o *TaskOptions) {
		o.EnvVars = ev
	}
}

func NewTaskOptions(opts ...TaskOption) TaskOptions {
	options := TaskOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
