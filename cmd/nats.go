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
	natsCmd = &cobra.Command{
		Use:   "nats",
		Short: "Manage the platform's nats broker",
		Long:  "Manage the platform's nats broker.",
	}

	validateNatsCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate nats",
		Long:  "Validate nats.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range service() {
				// get the steps
				steps, err := p.NatsSteps()
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

	planNatsCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plan nats",
		Long:  "Plan nats.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range nats() {
				// get the steps
				steps, err := p.NatsSteps()
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

	applyNatsCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply nats",
		Long:  "Apply nats.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range nats() {
				// get the steps
				steps, err := p.NatsSteps()
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

	destroyNatsCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy nats",
		Long:  "Destroy nats.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range nats() {
				// get the steps
				steps, err := p.NatsSteps()
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

func nats() []step.Platform {
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
	natsCmd.AddCommand(validateNatsCmd)
	natsCmd.AddCommand(planNatsCmd)
	natsCmd.AddCommand(applyNatsCmd)
	natsCmd.AddCommand(destroyNatsCmd)

	natsCmd.PersistentFlags().StringP("nats-namespace", "", "", "The namespace of nats")
	viper.BindPFlag("nats-namespace", natsCmd.PersistentFlags().Lookup("nats-namespace"))

	natsCmd.PersistentFlags().StringP("image-pull-policy", "", "", "The k8s image pull policy")
	viper.BindPFlag("image-pull-policy", natsCmd.PersistentFlags().Lookup("image-pull-policy"))

	rootCmd.AddCommand(natsCmd)
}
