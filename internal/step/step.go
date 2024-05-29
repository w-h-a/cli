package step

import "github.com/w-h-a/cli/internal/task"

type Step []task.Task

func ExecuteValidate(steps []Step) error {
	for _, step := range steps {
		for _, t := range step {
			defer t.Finalize()
			if err := t.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
