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
	cockroachCmd = &cobra.Command{
		Use:   "cockroach",
		Short: "Manage the platform's cockroach databases",
		Long:  "Manage the platform's cockroach databases.",
	}

	validateCockroachCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate cockroach",
		Long:  "Validate cockroach.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range service() {
				// get the steps
				steps, err := p.CockroachSteps()
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

	planCockroachCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plan cockroach",
		Long:  "Plan cockroach.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range cockroach() {
				// get the steps
				steps, err := p.CockroachSteps()
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

	applyCockroachCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply cockroach",
		Long:  "Apply cockroach.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range cockroach() {
				// get the steps
				steps, err := p.CockroachSteps()
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

	destroyCockroachCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy cockroach",
		Long:  "Destroy cockroach.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range cockroach() {
				// get the steps
				steps, err := p.CockroachSteps()
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

func cockroach() []step.Platform {
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
	cockroachCmd.AddCommand(validateCockroachCmd)
	cockroachCmd.AddCommand(planCockroachCmd)
	cockroachCmd.AddCommand(applyCockroachCmd)
	cockroachCmd.AddCommand(destroyCockroachCmd)

	cockroachCmd.PersistentFlags().StringP("cockroachdb-namespace", "", "", "The namespace of cockroachdb")
	viper.BindPFlag("cockroachdb-namespace", cockroachCmd.PersistentFlags().Lookup("cockroachdb-namespace"))

	cockroachCmd.PersistentFlags().StringP("image-pull-policy", "", "", "The k8s image pull policy")
	viper.BindPFlag("image-pull-policy", cockroachCmd.PersistentFlags().Lookup("image-pull-policy"))

	rootCmd.AddCommand(cockroachCmd)
}
