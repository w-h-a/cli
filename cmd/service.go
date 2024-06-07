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
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Manage the platform's service's",
		Long:  "Manage the platform's k8s cluster.",
	}

	validateServiceCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate service",
		Long:  "Validate service.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range service() {
				// get the steps
				steps, err := p.ServiceSteps()
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

	planServiceCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plan service",
		Long:  "Plan service.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range service() {
				// get the steps
				steps, err := p.ServiceSteps()
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

	applyServiceCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply service",
		Long:  "Apply service.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range service() {
				// get the steps
				steps, err := p.ServiceSteps()
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

	destroyServiceCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy service",
		Long:  "Destroy service.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range service() {
				// get the steps
				steps, err := p.ServiceSteps()
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

func service() []step.Platform {
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
	serviceCmd.AddCommand(validateServiceCmd)
	serviceCmd.AddCommand(planServiceCmd)
	serviceCmd.AddCommand(applyServiceCmd)
	serviceCmd.AddCommand(destroyServiceCmd)

	serviceCmd.PersistentFlags().StringP("service-namespace", "", "", "The service's namespace")
	viper.BindPFlag("service-namespace", serviceCmd.PersistentFlags().Lookup("service-namespace"))

	serviceCmd.PersistentFlags().StringP("service-name", "", "", "The service's name")
	viper.BindPFlag("service-name", serviceCmd.PersistentFlags().Lookup("service-name"))

	serviceCmd.PersistentFlags().StringP("service-version", "", "", "The service's version")
	viper.BindPFlag("service-version", serviceCmd.PersistentFlags().Lookup("service-version"))

	serviceCmd.PersistentFlags().StringP("service-port", "", "", "The service's port")
	viper.BindPFlag("service-port", serviceCmd.PersistentFlags().Lookup("service-port"))

	serviceCmd.PersistentFlags().StringP("service-image", "", "", "The service's base repo/image")
	viper.BindPFlag("service-image", serviceCmd.PersistentFlags().Lookup("service-image"))

	rootCmd.AddCommand(serviceCmd)
}
