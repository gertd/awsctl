package cmd

import (
	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop instances",
	RunE:  runStop,
}

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVar(&name, "name", "", "AWS instance name: default all")
}

func runStop(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region)

	client.StopInstance(name)
	client.GetInstances().Print()

	return nil
}
