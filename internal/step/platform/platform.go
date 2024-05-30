package platform

import (
	"fmt"

	"github.com/w-h-a/cli/internal/step"
	"github.com/w-h-a/cli/internal/step/kubernetes"
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
		task.TaskWithName(p.Name + "check my state"),
	)

	steps = append(steps, step.Step{stateChecker})

	for _, r := range p.Regions {
		// 2.1 create kubernetes
		k := &kubernetes.Kubernetes{
			Name:     p.Name,
			Env:      p.Env,
			Provider: r.Provider,
			Region:   r.Region,
		}

		cluster, err := k.Steps()
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes steps: %v", err)
		}

		steps = append(steps, cluster...)
	}

	return steps, nil
}
