package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/step"
	"github.com/w-h-a/cli/internal/step/platform"
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
)

func infra() []platform.Platform {
	// TODO: take this stuff as cli input
	platforms := []platform.Platform{
		{
			Name: "platform",
			Env:  "prod",
			Regions: []platform.Region{
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

	rootCmd.AddCommand(infraCmd)
}
