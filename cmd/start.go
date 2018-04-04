package cmd

import (
	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start instances",
	RunE:  runStart,
}

var (
	name string
)

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&name, "name", "", "AWS instance name: default all")
}

func runStart(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region)
	defer client.Close()

	client.StartInstance(name)
	client.GetInstances().Print()

	return nil
}
