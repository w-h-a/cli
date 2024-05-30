package kubernetes

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/step"
	"github.com/w-h-a/cli/internal/task"
	"github.com/w-h-a/cli/internal/task/terraform"
)

type Kubernetes struct {
	Name     string
	Env      string
	Provider string
	Region   string
}

func (k *Kubernetes) Steps() ([]step.Step, error) {
	steps := []step.Step{}

	k8sName := k.internalName("k8s")

	vars := map[string]string{}

	vars["token"] = viper.GetString("token")
	vars["name"] = k.Name
	vars["env"] = k.Env
	vars["region"] = k.Region

	k8s := terraform.NewTask(
		task.TaskWithName(k8sName),
		task.TaskWithSource(fmt.Sprintf("https://github.com/w-h-a/kubernetes-%s.git", k.Provider)),
		task.TaskWithPath(fmt.Sprintf("/tmp/%s", k8sName)),
		terraform.TerraformWithVars(vars),
	)

	steps = append(steps, step.Step{k8s}, step.Step{})

	return steps, nil
}

func (k *Kubernetes) internalName(name string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", k.Name, k.Env, k.Region, k.Provider, name)
}
