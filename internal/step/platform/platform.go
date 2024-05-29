package platform

import (
	"github.com/w-h-a/cli/internal/step"
	"github.com/w-h-a/cli/internal/task"
	"github.com/w-h-a/cli/internal/task/state"
)

type Platform struct {
	Name    string
	Env     string
	Domain  string
	Regions []Region
}

func (p *Platform) InfraSteps() ([]step.Step, error) {
	steps := []step.Step{}

	// 1: ensure remote state is available
	stateChecker := state.NewTask(
		task.TaskWithName(p.Name),
	)

	steps = append(steps, step.Step{stateChecker})

	return steps, nil
}
