package task

type TaskOption func(o *TaskOptions)

type TaskOptions struct {
	Name string
}

func TaskWithName(n string) TaskOption {
	return func(o *TaskOptions) {
		o.Name = n
	}
}

func NewTaskOptions(opts ...TaskOption) TaskOptions {
	options := TaskOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
