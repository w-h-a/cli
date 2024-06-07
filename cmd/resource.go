package cmd

import "github.com/spf13/cobra"

var (
	resourceCmd = &cobra.Command{
		Use:   "resource",
		Short: "Manage the platform's shared resources",
		Long:  "Manage the platform's shared resources.",
	}
)

func init() {
	rootCmd.AddCommand(resourceCmd)
}
