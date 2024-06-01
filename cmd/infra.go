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
	infraCmd = &cobra.Command{
		Use:   "infra",
		Short: "Manage the platform's infrastructure",
		Long:  "Manage the platform's infrastructure.",
	}

	validateInfraCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate terraform",
		Long:  "Validate terraform.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range infra() {
				// get the steps
				steps, err := p.InfraSteps()
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

	planInfraCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plan terraform",
		Long:  "Plan terraform.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range infra() {
				// get the steps
				steps, err := p.InfraSteps()
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

	applyInfraCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply terraform",
		Long:  "Apply terraform.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range infra() {
				// get the steps
				steps, err := p.InfraSteps()
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

	destroyInfraCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy terraform",
		Long:  "Destroy terraform.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range infra() {
				// get the steps
				steps, err := p.InfraSteps()
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

func infra() []step.Platform {
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
	infraCmd.AddCommand(validateInfraCmd)
	infraCmd.AddCommand(planInfraCmd)
	infraCmd.AddCommand(applyInfraCmd)
	infraCmd.AddCommand(destroyInfraCmd)

	rootCmd.AddCommand(infraCmd)
}
