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
	runtimeCmd = &cobra.Command{
		Use:   "runtime",
		Short: "Manage the platform's runtime",
		Long:  "Manage the platform's runtime.",
	}

	validateRuntimeCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate runtime",
		Long:  "Validate runtime.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range runtime() {
				// get the steps
				steps, err := p.RuntimeSteps()
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

	planRuntimeCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plan runtime",
		Long:  "Plan runtime.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range runtime() {
				// get the steps
				steps, err := p.RuntimeSteps()
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

	applyRuntimeCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply runtime",
		Long:  "Apply runtime.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range runtime() {
				// get the steps
				steps, err := p.RuntimeSteps()
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

	destroyRuntimeCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy runtime",
		Long:  "Destroy runtime.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range runtime() {
				// get the steps
				steps, err := p.RuntimeSteps()
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

func runtime() []step.Platform {
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
	runtimeCmd.AddCommand(validateRuntimeCmd)
	runtimeCmd.AddCommand(planRuntimeCmd)
	runtimeCmd.AddCommand(applyRuntimeCmd)
	runtimeCmd.AddCommand(destroyRuntimeCmd)

	runtimeCmd.PersistentFlags().StringP("resource-namespace", "", "", "The namespace of shared resources")
	viper.BindPFlag("resource-namespace", runtimeCmd.PersistentFlags().Lookup("resource-namespace"))

	runtimeCmd.PersistentFlags().StringP("app-namespace", "", "", "The namespace of the app")
	viper.BindPFlag("app-namespace", runtimeCmd.PersistentFlags().Lookup("app-namespace"))

	runtimeCmd.PersistentFlags().StringP("service-namespace", "", "", "The service's namespace")
	viper.BindPFlag("service-namespace", runtimeCmd.PersistentFlags().Lookup("service-namespace"))

	runtimeCmd.PersistentFlags().StringP("service-name", "", "", "The service's name")
	viper.BindPFlag("service-name", runtimeCmd.PersistentFlags().Lookup("service-name"))

	runtimeCmd.PersistentFlags().StringP("service-version", "", "", "The service's version")
	viper.BindPFlag("service-version", runtimeCmd.PersistentFlags().Lookup("service-version"))

	runtimeCmd.PersistentFlags().StringP("service-port", "", "", "The service's port")
	viper.BindPFlag("service-port", runtimeCmd.PersistentFlags().Lookup("service-port"))

	runtimeCmd.PersistentFlags().StringP("service-image", "", "", "The service's base repo/image")
	viper.BindPFlag("service-image", runtimeCmd.PersistentFlags().Lookup("service-image"))

	runtimeCmd.PersistentFlags().StringP("image-pull-policy", "", "", "The k8s image pull policy")
	viper.BindPFlag("image-pull-policy", runtimeCmd.PersistentFlags().Lookup("image-pull-policy"))

	rootCmd.AddCommand(runtimeCmd)
}
