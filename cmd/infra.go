package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/step"
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
	// TODO: take this stuff as cli input
	platforms := []step.Platform{
		{
			Name: "platform",
			Env:  "prod",
			Regions: []step.Region{
				{
					Provider: "do",
					Region:   viper.GetString("do-region"),
				},
			},
		},
	}

	return platforms
}

func init() {
	infraCmd.AddCommand(validateInfraCmd)
	infraCmd.AddCommand(planInfraCmd)
	infraCmd.AddCommand(applyInfraCmd)
	infraCmd.AddCommand(destroyInfraCmd)

	rootCmd.AddCommand(infraCmd)
}
