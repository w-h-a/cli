package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/step"
	"gopkg.in/yaml.v2"
)

var (
	k8sCmd = &cobra.Command{
		Use:   "k8s",
		Short: "Manage the platform's k8s cluster",
		Long:  "Manage the platform's k8s cluster.",
	}

	validateK8sCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate k8s",
		Long:  "Validate k8s.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range k8s() {
				// get the steps
				steps, err := p.K8sSteps()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}

				// validate them
				if err := step.ExecuteValidate(steps); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}
			}

			fmt.Println("validation succeeded")
		},
	}

	planK8sCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plan k8s",
		Long:  "Plan k8s.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range k8s() {
				// get the steps
				steps, err := p.K8sSteps()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}

				// plan them
				if err := step.ExecutePlan(steps); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}
			}

			fmt.Println("plan succeeded")
		},
	}

	applyK8sCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply k8s",
		Long:  "Apply k8s.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range k8s() {
				// get the steps
				steps, err := p.K8sSteps()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}

				// apply them
				if err := step.ExecuteApply(steps); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}
			}

			fmt.Println("apply succeeded")
		},
	}

	destroyK8sCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy k8s",
		Long:  "Destroy k8s.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range k8s() {
				// get the steps
				steps, err := p.K8sSteps()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}

				// destroy them
				if err := step.ExecuteDestroy(steps); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}
			}

			fmt.Println("destroy succeeded")
		},
	}
)

func k8s() []step.Platform {
	if len(viper.Get("config-file").(string)) == 0 {
		fmt.Fprintf(os.Stderr, "no platforms defined in the config file %s\n", viper.Get("config-file"))
		os.Exit(1)
	}

	configBytes, err := os.ReadFile(viper.Get("config-file").(string))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config file: %s\n", err.Error())
		os.Exit(1)
	}

	platforms := []step.Platform{}

	platform := step.Platform{}

	// TODO: figure out how to unmarshal array of platforms from yaml
	if err := yaml.Unmarshal(configBytes, &platform); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal config file: %s\n", err.Error())
		os.Exit(1)
	}

	platforms = append(platforms, platform)

	return platforms
}

func init() {
	k8sCmd.AddCommand(validateK8sCmd)
	k8sCmd.AddCommand(planK8sCmd)
	k8sCmd.AddCommand(applyK8sCmd)
	k8sCmd.AddCommand(destroyK8sCmd)

	rootCmd.AddCommand(k8sCmd)
}
