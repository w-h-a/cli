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
	Name    string   `yaml:"name"`
	Env     string   `yaml:"env"`
	Domain  string   `yaml:"domain,omitempty"`
	Regions []Region `yaml:"regions"`
}

type Region struct {
	Provider string `yaml:"provider"`
	Region   string `yaml:"region"`
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
		vars["app_namespace"] = strings.ToLower(fmt.Sprintf("%s-app", p.Name))
		vars["content_namespace"] = strings.ToLower(fmt.Sprintf("%s-content", p.Name))
		vars["misc_namespace"] = strings.ToLower(fmt.Sprintf("%s-misc", p.Name))

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

func (p *Platform) CockroachSteps() ([]Step, error) {
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

		// 2.2. cockroach
		cockroachName := p.internalName(r, fmt.Sprintf("%s.%s", "cockroachdb", viper.GetString("cockroachdb-namespace")))

		vars := map[string]string{}

		vars["cockroachdb_namespace"] = viper.GetString("cockroachdb-namespace")
		vars["image_pull_policy"] = viper.GetString("image-pull-policy")

		service := terraform.NewTask(
			task.TaskWithName(cockroachName),
			task.TaskWithSource(fmt.Sprintf("%s/kubernetes-cockroach.git", viper.GetString("base-source"))),
			task.TaskWithPath(fmt.Sprintf("/tmp/%s", cockroachName)),
			task.TaskWithEnvVars(env),
			terraform.TerraformWithVars(vars),
		)

		steps = append(steps, Step{service})
	}

	return steps, nil
}

func (p *Platform) NatsSteps() ([]Step, error) {
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

		// 2.2. nats
		natsName := p.internalName(r, fmt.Sprintf("%s.%s", "nats", viper.GetString("nats-namespace")))

		vars := map[string]string{}

		vars["nats_namespace"] = viper.GetString("nats-namespace")
		vars["image_pull_policy"] = viper.GetString("image-pull-policy")

		service := terraform.NewTask(
			task.TaskWithName(natsName),
			task.TaskWithSource(fmt.Sprintf("%s/kubernetes-nats.git", viper.GetString("base-source"))),
			task.TaskWithPath(fmt.Sprintf("/tmp/%s", natsName)),
			task.TaskWithEnvVars(env),
			terraform.TerraformWithVars(vars),
		)

		steps = append(steps, Step{service})
	}

	return steps, nil
}

func (p *Platform) RuntimeSteps() ([]Step, error) {
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

		// 2.2. runtime
		serviceName := p.internalName(r, fmt.Sprintf("%s.%s", viper.GetString("runtime-name"), viper.GetString("runtime-namespace")))

		vars := map[string]string{}

		vars["service_namespace"] = viper.GetString("runtime-namespace")
		vars["service_name"] = viper.GetString("runtime-name")
		vars["service_version"] = viper.GetString("runtime-version")
		vars["service_port"] = viper.GetString("runtime-port")
		vars["service_image"] = viper.GetString("runtime-image")
		vars["image_pull_policy"] = viper.GetString("runtime-pull-policy")

		service := terraform.NewTask(
			task.TaskWithName(serviceName),
			task.TaskWithSource(fmt.Sprintf("%s/kubernetes-runtime.git", viper.GetString("base-source"))),
			task.TaskWithPath(fmt.Sprintf("/tmp/%s", serviceName)),
			task.TaskWithEnvVars(env),
			terraform.TerraformWithVars(vars),
		)

		steps = append(steps, Step{service})
	}

	return steps, nil
}

func (p *Platform) ServiceSteps() ([]Step, error) {
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

		// 2.2. service
		serviceName := p.internalName(r, fmt.Sprintf("%s.%s", viper.GetString("service-name"), viper.GetString("service-namespace")))

		vars := map[string]string{}

		vars["resource_namespace"] = viper.GetString("resource-namespace")
		vars["app_namespace"] = viper.GetString("app-namespace")
		vars["service_namespace"] = viper.GetString("service-namespace")
		vars["service_name"] = viper.GetString("service-name")
		vars["service_version"] = viper.GetString("service-version")
		vars["service_type"] = viper.GetString("service-type")
		vars["service_port"] = viper.GetString("service-port")
		vars["node_port"] = viper.GetString("node-port")
		vars["service_image"] = viper.GetString("service-image")
		vars["image_pull_policy"] = viper.GetString("image-pull-policy")
		vars["admin"] = viper.GetString("admin")
		vars["secret"] = viper.GetString("secret")
		vars["payment_key"] = viper.GetString("payment-key")
		vars["enable_tls"] = viper.GetString("enable-tls")
		vars["cert_provider"] = viper.GetString("cert-provider")
		vars["hosts"] = viper.GetString("hosts")
		vars["aws_access_key"] = viper.GetString("aws-access-key")
		vars["aws_secret_access_key"] = viper.GetString("aws-secret-access-key")

		service := terraform.NewTask(
			task.TaskWithName(serviceName),
			task.TaskWithSource(fmt.Sprintf("%s/kubernetes-service.git", viper.GetString("base-source"))),
			task.TaskWithPath(fmt.Sprintf("/tmp/%s", serviceName)),
			task.TaskWithEnvVars(env),
			terraform.TerraformWithVars(vars),
		)

		steps = append(steps, Step{service})
	}

	return steps, nil
}

func (p *Platform) internalName(r Region, name string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", p.Name, p.Env, r.Region, r.Provider, name)
}
