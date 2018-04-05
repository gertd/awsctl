package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	region     string
	name       string
	keyFile    string
	instanceID string
)

// RootCmd -- represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "awsctl",
	Short: "AWS instance manager",
}

// Execute -- adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	RootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region: default all")
}
