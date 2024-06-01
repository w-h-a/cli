package step

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/task"
	"github.com/w-h-a/cli/internal/task/state"
	"github.com/w-h-a/cli/internal/task/terraform"
)

type Platform struct {
	Name    string
	Env     string
	Domain  string
	Regions []Region
}

type Region struct {
	Provider string
	Region   string
}

func (p *Platform) InfraSteps() ([]Step, error) {
	steps := []Step{}

	// 1. ensure remote state is available
	stateChecker := state.NewTask(
		task.TaskWithName(p.Name + "check my state"),
	)

	steps = append(steps, Step{stateChecker})

	for _, r := range p.Regions {
		// 2.1 kubernetes cluster
		k8sName := p.internalName(r, "k8s")

		vars := map[string]string{}

		vars["do_token"] = viper.GetString("do-token")
		vars["name"] = p.Name
		vars["region"] = r.Region

		k8s := terraform.NewTask(
			task.TaskWithName(k8sName),
			task.TaskWithSource(fmt.Sprintf("%s/kubernetes-%s.git", viper.GetString("base-source"), r.Provider)),
			task.TaskWithPath(fmt.Sprintf("/tmp/%s", k8sName)),
			terraform.TerraformWithVars(vars),
		)

		steps = append(steps, Step{k8s})
	}

	return steps, nil
}

func (p *Platform) K8sSteps() ([]Step, error) {
	steps := []Step{}

	// 1. ensure remote state is available
	stateChecker := state.NewTask(
		task.TaskWithName(p.Name + "check my state"),
	)

	steps = append(steps, Step{stateChecker})

	for _, r := range p.Regions {
		env := map[string]string{}

		env["KUBE_CONFIG_PATH"] = "~/.kube/config"

		// 2.1. kubeconfig
		if r.Provider != "kind" {
			configName := p.internalName(r, "kubeconfig")

			remoteStates := map[string]string{}

			remoteStates["k8s"] = p.internalName(r, "k8s")

			vars := map[string]string{}

			vars["do_token"] = viper.GetString("do-token")
			vars["kubernetes"] = r.Provider

			config := terraform.NewTask(
				task.TaskWithName(configName),
				task.TaskWithSource(fmt.Sprintf("%s/kubeconfig.git", viper.GetString("base-source"))),
				task.TaskWithPath(fmt.Sprintf("/tmp/%s", configName)),
				terraform.TerraformWithRemoteStates(remoteStates),
				terraform.TerraformWithVars(vars),
			)

			steps = append(steps, Step{config})

			env["KUBE_CONFIG_PATH"] = fmt.Sprintf("/tmp/%s/kubeconfig", configName)
		}

		// 2.2. namespaces
		namespaceName := p.internalName(r, "namespaces")

		vars := map[string]string{}

		vars["resource_namespace"] = strings.ToLower(fmt.Sprintf("%s-resource", p.Name))

		namespace := terraform.NewTask(
			task.TaskWithName(namespaceName),
			task.TaskWithSource(fmt.Sprintf("%s/kubernetes-namespaces.git", viper.GetString("base-source"))),
			task.TaskWithPath(fmt.Sprintf("/tmp/%s", namespaceName)),
			task.TaskWithEnvVars(env),
			terraform.TerraformWithVars(vars),
		)

		steps = append(steps, Step{namespace})
	}

	return steps, nil
}

func (p *Platform) internalName(r Region, name string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", p.Name, p.Env, r.Region, r.Provider, name)
}
