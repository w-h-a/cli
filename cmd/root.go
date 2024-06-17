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
	Long:  `The cli. Deploy base infra, shared ad hoc resources, and services.`,
}

func viperConfig() {
	viper.SetDefault("state-store", "aws")
	viper.SetDefault("aws-region", "us-west-2")
	viper.SetDefault("aws-s3-bucket", "wha-infra-terraform-state")
	viper.SetDefault("aws-dynamodb-table", "wha-infra-terraform-lock")
	viper.SetDefault("base-source", "https://github.com/w-h-a")
	viper.SetDefault("node-port", "0")
}

func init() {
	cobra.OnInitialize(viperConfig)

	rootCmd.PersistentFlags().StringP("do-token", "d", "", "DO provider token")
	viper.BindPFlag("do-token", rootCmd.PersistentFlags().Lookup("do-token"))

	rootCmd.PersistentFlags().StringP("aws-s3-bucket", "b", "", "AWS S3 bucket name")
	viper.BindPFlag("aws-s3-bucket", rootCmd.PersistentFlags().Lookup("aws-s3-bucket"))

	rootCmd.PersistentFlags().StringP("aws-dynamodb-table", "t", "", "AWS Dynamodb Table")
	viper.BindPFlag("aws-dynamodb-table", rootCmd.PersistentFlags().Lookup("aws-dynamodb-table"))

	rootCmd.PersistentFlags().StringP("base-source", "s", "", "Base source")
	viper.BindPFlag("base-source", rootCmd.PersistentFlags().Lookup("base-source"))

	rootCmd.PersistentFlags().StringP("config-file", "c", "", "Path to config file")
	viper.BindPFlag("config-file", rootCmd.PersistentFlags().Lookup("config-file"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
