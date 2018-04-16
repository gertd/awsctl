package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	defName       = ""
	defInstanceID = ""
	defRegion     = ""
	defKeyFile    = ""
	defAll        = false
)

var (
	region     string
	name       string
	keyFile    string
	instanceID string
	all        bool
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
	cobra.EnableCommandSorting = false

	RootCmd.PersistentFlags().StringVar(&region, "region", defRegion, "AWS region: default all")
}
