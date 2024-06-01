package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "The cli",
	Long:  `The cli. Deploy base infra, shared resources, and services.`,
}

func viperConfig() {
	viper.SetDefault("do-region", "sfo2")

	viper.SetDefault("state-store", "aws")
	viper.SetDefault("aws-region", "us-west-2")

	// TODO: take via flags
	viper.SetDefault("aws-s3-bucket", "wha-infra-terraform-state")
	viper.SetDefault("aws-dynamodb-table", "wha-infra-terraform-lock")
	viper.SetDefault("base-source", "https://github.com/w-h-a")
}

func init() {
	cobra.OnInitialize(viperConfig)

	rootCmd.PersistentFlags().StringP("do-token", "d", "", "DO provider token")
	viper.BindPFlag("do-token", rootCmd.PersistentFlags().Lookup("do-token"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
