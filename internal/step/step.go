package step

import (
	"strings"

	"github.com/w-h-a/cli/internal/task"
)

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

func ExecutePlan(steps []Step) error {
	for _, step := range steps {
		for _, t := range step {
			defer t.Finalize()

			if err := t.Validate(); err != nil {
				return err
			}

			if err := t.Plan(); err != nil {
				return err
			}
		}
	}

	return nil
}

func ExecuteApply(steps []Step) error {
	for _, step := range steps {
		for _, t := range step {
			defer t.Finalize()

			if err := t.Validate(); err != nil {
				return err
			}

			if err := t.Apply(); err != nil {
				return err
			}
		}
	}

	return nil
}

func ExecuteDestroy(steps []Step) error {
	// first find the kubeconfig, apply it, and then defer its destruction
	for _, step := range steps {
		for _, t := range step {
			if strings.Contains(t.Options().Source, "kubeconfig") {
				defer t.Finalize()

				if err := t.Validate(); err != nil {
					return err
				}

				if err := t.Apply(); err != nil {
					return err
				}

				defer t.Destroy()
			}
		}
	}

	for i := len(steps) - 1; i >= 0; i-- {
		for _, t := range steps[i] {
			if !strings.Contains(t.Options().Source, "kubeconfig") {
				defer t.Finalize()

				if err := t.Validate(); err != nil {
					return err
				}

				if err := t.Destroy(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
